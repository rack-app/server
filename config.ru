
require 'rack/app'

app = lambda do |_env|
  resp = Rack::Response.new
  resp.write('OK')
  resp.finish
end

require 'socket'
require 'timeout'
require 'rack/content_length'
require 'rack/rewindable_input'

module Rack::App::Handler
  ::Rack::Handler.register 'rack-app', 'Rack::App::Handler'

  extend(self)

  class Server
    def initialize(app, socket_path)
      @socket_path = socket_path
      @mutex = Mutex.new
      @app = app
    end

    def start
      @server = UNIXServer.new(@socket_path)
      loop { handle_request(@server.accept) }
    rescue IOError, Errno::EBADF
      File.delete(@socket_path) if File.exist?(@socket_path)
    end

    def stop
      Thread.new do
        @mutex.synchronize do
          @server.close
        end
      end
    end

    private

    def handle_request(socket)
      @mutex.synchronize do
        env = receive_env(socket)
        resp = response_for(env)
        handle_response(socket, resp)
        socket.close
      end
    end

    def handle_response(socket, raw_rack_resp)
      socket.puts(JSON.dump(resp_conf_by(raw_rack_resp)))

      raw_rack_resp[2].each do |chunk|
        socket.print(chunk)
      end

      socket.flush
    end

    def resp_conf_by(rack_resp)
      {
        'status' => rack_resp[0],
        'headers' => rack_resp[1],
        'length' => (rack_resp[1][::Rack::CONTENT_LENGTH] || -1).to_i
      }
    end

    def response_for(env)
      @app.call(env)
    rescue Exception => ex
      env[::Rack::RACK_ERRORS] << ex.message << "\n"
      [500, {}, []]
    end

    require 'json'
    def receive_env(socket)
      env_base = JSON.parse(socket.gets)
      env_base.merge(::Rack::RACK_VERSION => Rack::VERSION,
                     ::Rack::RACK_INPUT        => Rack::RewindableInput.new(socket),
                     ::Rack::RACK_ERRORS       => $stderr,
                     ::Rack::RACK_MULTITHREAD  => false,
                     ::Rack::RACK_MULTIPROCESS => false,
                     ::Rack::RACK_RUNONCE      => false,
                     ::Rack::RACK_URL_SCHEME   => env_base['SCHEME'].downcase)
    end
  end

  def run(app, options)
    socket_path = options[:SP] || abort('Missing Socket file path')
    server = self::Server.new(app, socket_path)
    trap_signals_for(server)
    server.start
  end

  def valid_options
    { 'SP=[SOCKET_FILE_PATH]' => 'socket file path where the socket should be created' }
  end

  private

  def trap_signals_for(receiver)
    %w[INT TERM].each { |sig| ::Signal.trap(sig) { receiver.stop } }
  end
end

run app
