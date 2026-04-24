#include <direct.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/**
 * This program prints the current working directory.
 * It accepts an optional '-u' flag to output the path in Unix style (using '/' as separators).
 * 
 * Usage: wpwd [-u]
 * 
 * @param argc The number of command-line arguments.
 * @param argv The array of command-line arguments.
 * @return 0 on success, non-zero on failure.
 */
int main(int argc, char *argv[]) {
	int unix_style = 0;

	if (argc >= 2) {
        if (strcmp(argv[1], "-u") == 0) {
			unix_style = 1;
		} else {
		    fprintf(stderr, "Usage: wpwd [-u]\n");
		    return 1;
        }
	}

	char *cwd = _getcwd(NULL, 0);
	if (cwd == NULL) {
		perror("pwd");
		return 1;
	}

	if (!unix_style) {
        // Output the path as-is (Windows style)
		printf("%s\n", cwd);
	} else {
        // Convert backslashes to forward slashes for Unix style output
        printf("/");
		for (char *p = cwd; *p != '\0'; p++) {
			putchar(*p == '\\' ? '/' : *p);
		}
		putchar('\n');
	}

	free(cwd);
	return 0;
}
