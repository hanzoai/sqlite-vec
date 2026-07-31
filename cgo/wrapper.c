/* The one compilation unit for sqlite-vec in a Go build. sqlite-vec.c lives at
 * the repo root (with the -rescore/-diskann fragments it #include's); compiling
 * it via this wrapper keeps the root a plain C source tree that future upstream
 * merges land on untouched. */
#include "../sqlite-vec.c"
