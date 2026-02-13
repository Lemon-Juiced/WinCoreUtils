#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <windows.h>

/**
 * This is a simple wrapper around wls.exe that allows it to be invoked as "wla" or "wla.exe" to enable the long listing format by default with hidden files.
 * 
 * Usage:
 *  wla [directory]
 * 
 * @author Lemon
 * @param argc The number of command-line arguments.
 * @param argv The array of command-line arguments.
 * @return The exit code of the wls.exe process, or 1 if process creation fails.
 */
int main(int argc, char **argv) {
    int bufsize = 1024;
    for (int i = 1; i < argc; i++) bufsize += (int)strlen(argv[i]) + 3;
    char *cmd = (char*)malloc(bufsize);
    if (!cmd) return 2;

    strcpy(cmd, "wls.exe -l -a");
    for (int i = 1; i < argc; i++) {
        strcat(cmd, " ");
        int needq = strchr(argv[i], ' ') != NULL;
        if (needq) strcat(cmd, "\"");
        strcat(cmd, argv[i]);
        if (needq) strcat(cmd, "\"");
    }

    STARTUPINFOA si;
    PROCESS_INFORMATION pi;
    ZeroMemory(&si, sizeof(si));
    si.cb = sizeof(si);
    ZeroMemory(&pi, sizeof(pi));

    if (!CreateProcessA(NULL, cmd, NULL, NULL, FALSE, 0, NULL, NULL, &si, &pi)) {
        fprintf(stderr, "CreateProcess failed: %lu\n", GetLastError());
        free(cmd);
        return 1;
    }

    WaitForSingleObject(pi.hProcess, INFINITE);
    DWORD exitCode = 0;
    GetExitCodeProcess(pi.hProcess, &exitCode);
    CloseHandle(pi.hProcess);
    CloseHandle(pi.hThread);
    free(cmd);
    return (int) exitCode;
}
