package quadratic

import (
    "fmt"
    "math"
)

// Globals
var DEBUG = false


// Notes:
// Takes coefficients a, b, c as per ax^2 + bx + c = 0.
// Outputs the roots.
func Quadratic(coeff_a, coeff_b, coeff_c float64) (root1, root2 complex128) {

    // Check the discriminant sign to see if we have real or imaginary roots
    discriminant := math.Pow(coeff_b, 2) - (4 * coeff_a * coeff_c)

    // Check that a != 0, if a == 0, then return degenerative roots.
    if coeff_a != 0 {
        if discriminant >=0 { // No imaginary part, real roots
            root1 = complex((-1*coeff_b + math.Sqrt(discriminant))                                                                      / (2 * coeff_a),
                            0)
            root2 = complex((-1*coeff_b - math.Sqrt(discriminant))                                                                      / (2 * coeff_a),
                            0)
        } else { // Imaginary roots
            root1 = complex((-1*coeff_b) / (2 * coeff_a),
                            (math.Sqrt(-1 * discriminant)) / (2 * coeff_a))
            root2 = complex((-1*coeff_b) / (2 * coeff_a),
                            -1 * (math.Sqrt(-1 * discriminant)) / (2 * coeff_a))
        }
    } else {  // Degenerative roots, just a line not a quadratic.
        root1 = complex((-1 * coeff_c) / coeff_b, 0)
        root2 = root1
    }

    return root1, root2
}

func GenerateHtml(coeff_a, coeff_b, coeff_c float64, root1, root2 complex128)                    (htmlStr string) {
    // Print the coefficients with "x"'s.
    // The print the two roots underneath

    // Deal with special cases of a=0 (not quadratic) and a=1 (omit the 1).
    // Deal with the special cases of b=0 (no b term) and b=1 (omit the 1).
    // Deal with the special case of c=0 (no c term).
    // Display -ve b and c coeffs without the additional "+" sign.
    // Print complex numbers with no imaginary part as real numbers.

    if coeff_a == 0 {
        htmlStr += fmt.Sprintf("")
    } else if coeff_a == 1 {
        htmlStr += fmt.Sprintf("x^2")
    } else {
        htmlStr += fmt.Sprintf("%vx^2", coeff_a)
    }

    if coeff_b < 0 {
        htmlStr += fmt.Sprintf(" - %vx", math.Abs(coeff_b))
    } else {
        if coeff_b == 0 {
            htmlStr += fmt.Sprintf("")
        } else if coeff_b == 1 {
            htmlStr += fmt.Sprintf(" + x")
        } else {
            htmlStr += fmt.Sprintf(" + %vx", coeff_b)
        }
    }

    if coeff_c < 0 {
        htmlStr += fmt.Sprintf(" - %v", math.Abs(coeff_c))
    } else {
        if coeff_c == 0 {
            htmlStr += fmt.Sprintf("")
        } else {
            htmlStr += fmt.Sprintf(" + %v", coeff_c)
        }
    }

    htmlStr += fmt.Sprintf(" = 0<br>")

    if imag(root1) == 0 {
        htmlStr += fmt.Sprintf("root1 = %.3g", real(root1))
    } else {
        htmlStr += fmt.Sprintf("root1 = %.3g", root1)
    }
    if imag(root2) == 0 {
        htmlStr += fmt.Sprintf(", root2 = %.3g", real(root2))
    } else {
        htmlStr += fmt.Sprintf(", root2 = %.3g", root2)
    }

    return htmlStr
}

func SelfTest() (bool) {

    result := true

    result = result && TestSingleQuadratic(1.0, -5.0, 6.0,
                                     complex(3, 0), complex(2,0))

    result = result && TestSingleQuadratic(1.0, 3.0, 2.5,
                                     complex(-1.5, 0.5), complex(-1.5,-0.5))

    result = result && TestSingleQuadratic(2.0, 1.0, -1.0,
                                     complex(0.5, 0), complex(-1,0))

    result = result && TestSingleQuadratic(1.0, 0.0, -9.0,
                                     complex(3, 0), complex(-3,0))

    result = result && TestSingleQuadratic(1.0, -4.0, 0.0,
                                     complex(4, 0), complex(0,0))

    if result {
        fmt.Printf("... Passed.\n")
    } else {
        fmt.Printf("... FAILED.\n")
    }

    return result
}

func TestSingleQuadratic(coeff_a, coeff_b, coeff_c float64,
                   expRoot1, expRoot2 complex128) (result bool) {
    result = true

    root1, root2 := Quadratic(coeff_a, coeff_b, coeff_c)

    result = root1 == expRoot1 && root2 == expRoot2
    if !result {
        fmt.Printf("root1: expected %.3g, got %.3g\n", expRoot1, root1)
        fmt.Printf("root2: expected %.3g, got %.3g\n", expRoot2, root2)
    } else {
        fmt.Printf("%v\n", GenerateHtml(coeff_a, coeff_b, coeff_c,
                                        root1, root2))
    }

    return result
}
