
using Service;
using System;
using System.Collections.Generic;
using System.IO;
using System.Net.Sockets;
using System.Threading.Tasks;

namespace Service
{
    public class SocketManager
    {
        private Socket socket;
        private Stream? stream;

        public SocketManager(string ipAddress, int port)
        {
            socket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            try
            {
                socket.Connect(ipAddress, port);
                stream = new NetworkStream(socket);
            }
            catch
            {
                Methods.Error("Connecting to the server!");
            }
        }
        
        ~SocketManager()
        {
            try
            {
                if (socket != null && socket.Connected)
                    socket.Shutdown(SocketShutdown.Both);
            }
            catch (Exception ex)
            {
                Methods.Error($"Error disconnecting client: {ex.Message}");
            }
            finally
            {
                socket?.Close();
            }
        }

        private async Task<byte[]> _recieveData()
        {
            byte[] _bytes = new byte[1024];

            int bytesRead = await stream!.ReadAsync(_bytes);
            byte[] receivedData = new byte[bytesRead];
            Array.Copy(_bytes, receivedData, bytesRead);

            return receivedData;
        }


        public async Task SendData(byte[] data)
        {
            try
            {
                await stream!.WriteAsync(data);
            }
            catch
            {
                Methods.Error("sending data to server! (it might've closed)");
            }
        }

        public async Task<Data?> ReceiveData()
        {
            try
            {
                byte[] bytesRead = await _recieveData();
                return new Data(bytesRead);
            }
            catch { 
                Methods.Error("receiving data from server! (it might've closed)");
                return null;
            }
        }
    }
}