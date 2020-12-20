Euclid's algorithm

https://en.wikipedia.org/wiki/Euclidean_algorithm#Procedure

Implementations of the algorithm may be expressed in pseudocode. For example, the division-based version may be programmed as[19]

function gcd(a, b)
    while b ≠ 0
        t := b
        b := a mod b
        a := t
    return a
At the beginning of the kth iteration, the variable b holds the latest remainder rk−1, whereas the variable a holds its predecessor, rk−2. The step b := a mod b is equivalent to the above recursion formula rk ≡ rk−2 mod rk−1. The temporary variable t holds the value of rk−1 while the next remainder rk is being calculated. At the end of the loop iteration, the variable b holds the remainder rk, whereas the variable a holds its predecessor, rk−1.
