using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Service
{
    using System;
    using System.Diagnostics;
    using System.IO;
    using System.Linq;
    using Microsoft.Win32;

    class ProgramStarter
    {
        public static void StartProgram(string programName)
        {
            // Get the user and system PATH variables and merge them.
            string[] pathDirectories = GetMergedPathDirectories();

            // Try to find the program in the merged PATH directories.
            string programPath = FindProgramInPath(programName, pathDirectories);

            if (!string.IsNullOrEmpty(programPath))
            {
                Process.Start(programPath);
                Console.WriteLine($"Started program: {programPath}");
            }
            else
            {
                Console.WriteLine($"Program '{programName}' not found in the PATH.");
            }
        }

        private static string[] GetMergedPathDirectories()
        {
            string userPath = Environment.GetEnvironmentVariable("PATH", EnvironmentVariableTarget.User);
            string systemPath = Environment.GetEnvironmentVariable("PATH", EnvironmentVariableTarget.Machine);

            if (string.IsNullOrEmpty(userPath) && string.IsNullOrEmpty(systemPath))
            {
                return new string[0]; // Neither user nor system PATH set.
            }

            if (string.IsNullOrEmpty(userPath))
            {
                return systemPath.Split(';'); // Only system PATH set.
            }

            if (string.IsNullOrEmpty(systemPath))
            {
                return userPath.Split(';'); // Only user PATH set.
            }

            // Merge user and system PATH variables, removing duplicates.
            var mergedPaths = new HashSet<string>(userPath.Split(';'), StringComparer.OrdinalIgnoreCase);
            foreach (var path in systemPath.Split(';'))
            {
                mergedPaths.Add(path);
            }

            return mergedPaths.ToArray();
        }

        private static string FindProgramInPath(string programName, string[] pathDirectories)
        {
            foreach (var directory in pathDirectories)
            {
                var fullPath = Path.Combine(directory, programName);
                if (File.Exists(fullPath))
                {
                    return fullPath;
                }
            }
            return null;
        }
    }
}
