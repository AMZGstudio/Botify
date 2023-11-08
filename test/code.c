#include <windows.h>
#include <stdio.h>

int main()
{
    // Specify the path to Notepad.exe
    LPCTSTR applicationName;
    applicationName = "C:\\Windows\\System32\\notepad.exe";

    // Initialize the STARTUPINFO and PROCESS_INFORMATION structures
    STARTUPINFO startupInfo;
    PROCESS_INFORMATION processInfo;
    ZeroMemory(&startupInfo, sizeof(STARTUPINFO));
    startupInfo.cb = sizeof(STARTUPINFO);

    // Create the Notepad process
    CreateProcess(
        applicationName,
        NULL,
        NULL,
        NULL,
        FALSE,
        0,
        NULL,
        NULL,
        &startupInfo,
        &processInfo
    );

    return 0;
}
