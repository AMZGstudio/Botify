using System;
using System.Collections.Generic;
using System.Text;
using System.Text.Json;
using System.Windows;

namespace Service
{
    public enum RequestType { LOGIN, SIGNUP};
    public class Methods
    {

        public static void closeProgram()
        {
            Environment.Exit(0);
        }
        public static void Error(string message)
        {
            Console.WriteLine("Error: " + message);
        }

        public static byte[] SerializeMessage(RequestType header, Dictionary<string, object> dict)
        {
            // Convert the header to bytes
            byte[] headerBytes = BitConverter.GetBytes((ushort)header);
            if (BitConverter.IsLittleEndian)
            {
                Array.Reverse(headerBytes);
            }

            // Convert the dictionary to a JSON object
            string jsonObject = JsonSerializer.Serialize(dict);

            // Convert the JSON object to bytes
            byte[] jsonObjectBytes = Encoding.UTF8.GetBytes(jsonObject);

            // Get the length of the JSON object
            int jsonObjectLength = jsonObjectBytes.Length;

            // Convert the length to bytes
            byte[] jsonObjectLengthBytes = BitConverter.GetBytes(jsonObjectLength);
            if (BitConverter.IsLittleEndian)
            {
                Array.Reverse(jsonObjectLengthBytes);
            }

            // Create the message
            byte[] message = new byte[headerBytes.Length + jsonObjectLengthBytes.Length + jsonObjectBytes.Length];
            Array.Copy(headerBytes, message, headerBytes.Length);
            Array.Copy(jsonObjectLengthBytes, 0, message, headerBytes.Length, jsonObjectLengthBytes.Length);
            Array.Copy(jsonObjectBytes, 0, message, headerBytes.Length + jsonObjectLengthBytes.Length, jsonObjectBytes.Length);

            return message;
        }
    }
}
