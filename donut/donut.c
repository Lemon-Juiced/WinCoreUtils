#include <stdio.h>
#include <stdlib.h>

/**
 * Checks if the point (dx, dy) is inside the circle defined by the outer radius squared.
 * 
 * @param dx The x-coordinate of the point relative to the center of the circle.
 * @param dy The y-coordinate of the point relative to the center of the circle.
 * @param outer_radius_sq The square of the outer radius of the circle.
 * @return 1 if the point is inside the circle, 0 otherwise.
 */
static int is_inside_circle(int dx, int dy, double outer_radius_sq) {
    return (double)(dx * dx + dy * dy) <= outer_radius_sq;
}

/**
 * Prints an error message to stderr in red and exits the program with a non-zero status code.
 * 
 * @param message The error message to be printed.
 */
static void print_error_and_exit(const char *message) {
    fprintf(stderr, "\033[1;31m%s\033[0m\n", message);
    exit(1);
}

/**
 * This program takes an odd integer greater than 1 and creates the shape of a * donut using '#' and '-' characters.
 * The '#' represents the perimeter of the donut, while the '-' represents the inner area of the donut.
 * The integer represents the width of the donut at its widest point.
 * 
 * Even sizes are not accepted because this rasterized circle (the "donut") requires a single center column and row.
 * 
 * Example (Input: 3):
 *   ###
 *  #---#
 * #-----#
 * #-----#
 * #-----#
 *  #---#
 *   ###
 * 
 * Usage: ./donut <odd size greater than 1>
 * 
 * @param argc The number of command-line arguments.
 * @param argv The array of command-line arguments.
 * @return 0 on success, non-zero on failure.
 */
int main(int argc, char *argv[]) {
    if (argc != 2) print_error_and_exit("Usage: ./donut <odd size greater than 1>");

    int size = atoi(argv[1]);
    if (size <= 1 || size % 2 == 0) print_error_and_exit("Error: Size must be an odd integer greater than 1.");

    int diameter = size * 2 + 1;
    int center = size;
    double outer_radius = size + 0.5;
    double outer_radius_sq = outer_radius * outer_radius;

    for (int i = 0; i < diameter; i++) {
        int dy = i - center;

        for (int j = 0; j < diameter; j++) {
            int dx = j - center;
            int is_inside = is_inside_circle(dx, dy, outer_radius_sq);
            int touches_outside;

            if (!is_inside) {
                printf(" ");
                continue;
            }

            touches_outside =
                !is_inside_circle(dx - 1, dy, outer_radius_sq) ||
                !is_inside_circle(dx + 1, dy, outer_radius_sq) ||
                !is_inside_circle(dx, dy - 1, outer_radius_sq) ||
                !is_inside_circle(dx, dy + 1, outer_radius_sq);

            if (touches_outside) {
                printf("#");
            } else {
                printf("-");
            }
        }
        printf("\n");
    }
    
    return 0;
}