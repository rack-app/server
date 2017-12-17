require 'socket'

socket_path = '/tmp/t.socket'

UNIXServer.open(socket_path) do |serv|
  UNIXSocket.open(socket_path) do |c|
    s = serv.accept

    c.send_io STDOUT
    stdout = s.recv_io

    p STDOUT.fileno #=> 1
    p stdout.fileno #=> 7

    stdout.puts 'hello' # outputs "hello\n" to standard output.
  end
end

File.delete(socket_path)
