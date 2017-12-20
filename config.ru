
require 'rack/app'

app = lambda do |_env|
  resp = Rack::Response.new
  resp.write('OK')
  resp.finish
end

require 'csv'
require 'socket'
require 'timeout'
require 'rack/content_length'
require 'rack/rewindable_input'

module Rack::App::Handler
  ::Rack::Handler.register 'rack-app-receiver', 'Rack::App::Handler'

  require 'socket'

  class Server
    def initialize(app, port)
      @port = port.to_i
      @app = app
      @wip = {}
    end

    def start
      @server = ::TCPServer.new(@port)
      @wip.clear

      loop do
        thr = Thread.start(@server.accept) do |s|
          handle(s)
          @wip.delete(Thread.current.__id__)
        end
        @wip[thr.__id__] = nil
      end
    rescue IOError, Errno::EBADF
      @server.close unless @server.closed?
    end

    def stop
      sleep(0.1) until @wip.empty?
      @server && @server.close
    end

    private

    def handle(socket)
      env = receive_env(socket)
      resp = response_for(env)

      send_headers(socket, resp[0], resp[1])
      send_body(socket, resp[2])
    ensure
      socket.flush
      socket.close
    end

    def send_headers(socket, code, headers)
      csv = CSV.new(socket, :col_sep => "\t")

      csv << [code]
      headers.each do |key, value|
        csv << [key, value]
      end

      socket.puts
    ensure
      socket.flush
    end

    def send_body(socket, eachable)
      eachable.each do |chunk|
        c = socket.write(chunk)
        puts("<- #{c}") if ENV["RACK_APP_DEBUG"]
      end
    ensure
      socket.flush
    end

    def response_for(env)
      @app.call(env)
    rescue Exception => ex
      env[::Rack::RACK_ERRORS] << ex.message << "\n"
      [500, {}, []]
    end

    def receive_env(socket)
      env = {}

      loop do
        line = socket.gets
        break if line.strip == ''
        row = CSV.parse(line, :col_sep => "\t", headers: %i[key value])
        env[row[:key].first] = row[:value].first.to_s
      end

      env.merge!(::Rack::RACK_VERSION      => Rack::VERSION,
                 ::Rack::RACK_INPUT        => Rack::RewindableInput.new(socket),
                 ::Rack::RACK_ERRORS       => $stderr,
                 ::Rack::RACK_MULTITHREAD  => true,
                 ::Rack::RACK_MULTIPROCESS => true,
                 ::Rack::RACK_RUNONCE      => false,
                 ::Rack::RACK_URL_SCHEME   => env['SCHEME'].downcase)
    end
  end

  module_function

  def run(app, options)
    server = self::Server.new(app, options[:Port] || '9292')
    %w[INT TERM].each { |sig| ::Signal.trap(sig) { server.stop } }
    server.start
  end

  def valid_options
    {}
  end
end

run app
