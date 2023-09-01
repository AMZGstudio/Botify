using System;
using System.Text;
using System.Text.Json;

namespace Service
{
    public enum RequestCode { Action, Query };

    public class Data
    {
        public RequestCode code { get; private set; }
        public string textData { get; private set; }

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

        public JsonDocument GetMessageData()
        {
            // Parse the JSON message and return a JsonDocument
            return JsonDocument.Parse(textData);
        }

        public string toString() { return "RequestCode: " + code + ", textData: " + textData; }
    }
}
