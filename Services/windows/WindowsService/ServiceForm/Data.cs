using System;
using System.Text;
using System.Text.Json;

namespace Service
{
    public enum RequestCode { Action, Query, QueryResponse, Trigger };

    public class Data
    {
        public RequestCode code { get; private set; }
        public string textData { get; private set; }

        public Data(RequestCode code, string textData)
        {
            this.code = code;
            this.textData = textData;
        }

        public Data(byte[] bytes)
        {
            if (bytes == null || bytes.Length < 7)
            {
                throw new ArgumentException("Invalid byte array. It must have at least 7 bytes.");
            }

            code = (RequestCode)bytes[0];

            // Extract the length of JSON data from the byte array
            int dataLength = BitConverter.ToInt32(bytes, 1);

            if (bytes.Length > 5)
            {
                textData = Encoding.UTF8.GetString(bytes, 6, bytes.Length-6);
            }
            else
            {
                textData = string.Empty; // No text data provided.
            }
        }

        public byte[] toBytes()
        {
            List<byte> bytes = new List<byte>();

            // Convert RequestCode enum to bytes
            bytes.Add((byte)code);

            // Convert textData to bytes
            byte[] textDataBytes = Encoding.UTF8.GetBytes(textData);

            // Add the length of textDataBytes as a 4-byte integer
            bytes.AddRange(BitConverter.GetBytes(textDataBytes.Length));

            // Add the textData bytes
            bytes.AddRange(textDataBytes);

            return bytes.ToArray();
        }

        public JsonDocument GetMessageData()
        {
            // Parse the JSON message and return a JsonDocument
            return JsonDocument.Parse(textData);
        }

        public string toString() { return "RequestCode: " + code + ", textData: " + textData; }
    }
}
