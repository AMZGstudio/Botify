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
            if (bytes == null || bytes.Length < 1)
            {
                throw new ArgumentException("Invalid byte array. It must have at least one byte.");
            }

            code = (RequestCode)bytes[0];

            if (bytes.Length > 1)
            {
                textData = Encoding.ASCII.GetString(bytes, 1, bytes.Length - 1);
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

    }
}
