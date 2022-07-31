package quadratic

import "testing"

func TestQuadratic(t *testing.T) {

    if !SelfTest() {
        t.Errorf("Quadratic self test failed")
    }
}

func TestGenerateHtml(t *testing.T) {

    htmlStrExp := "x^2 - 5x + 6 = 0<br>" +
                  "root1 = 3, root2 = 2"
    htmlStr := GenerateHtml(1, -5, 6, 3, 2)
    if htmlStr != htmlStrExp {
        t.Errorf("Expected html=\n\"%v\"\nbut got\n\"%v\"", htmlStrExp, htmlStr)
    }

    htmlStrExp = "x^2 + 3x + 2.5 = 0<br>" +
                 "root1 = (-1.5+0.5i), root2 = (-1.5-0.5i)"
    htmlStr = GenerateHtml(1, 3, 2.5, -1.5 + 0.5i, -1.5 - 0.5i)
    if htmlStr != htmlStrExp {
        t.Errorf("Expected html=\n\"%v\"\nbut got\n\"%v\"", htmlStrExp, htmlStr)
    }

    htmlStrExp = "2x^2 + x - 1 = 0<br>" +
                 "root1 = 0.5, root2 = -1"
    htmlStr = GenerateHtml(2, 1, -1, 0.5, -1)
    if htmlStr != htmlStrExp {
        t.Errorf("Expected html=\n\"%v\"\nbut got\n\"%v\"", htmlStrExp, htmlStr)
    }

    htmlStrExp = "x^2 - 9 = 0<br>" +
                 "root1 = 3, root2 = -3"
    htmlStr = GenerateHtml(1, 0, -9, 3, -3)
    if htmlStr != htmlStrExp {
        t.Errorf("Expected html=\n\"%v\"\nbut got\n\"%v\"", htmlStrExp, htmlStr)
    }

    htmlStrExp = "x^2 - 4x = 0<br>" +
                 "root1 = 4, root2 = 0"
    htmlStr = GenerateHtml(1, -4, 0, 4, 0)
    if htmlStr != htmlStrExp {
        t.Errorf("Expected html=\n\"%v\"\nbut got\n\"%v\"", htmlStrExp, htmlStr)
    }

    htmlStrExp = "x^2 + x - 1 = 0<br>" +
                 "root1 = 0.618, root2 = -1.62"
    htmlStr = GenerateHtml(1, 1, -1, 0.618033988, -1.61803399)
    if htmlStr != htmlStrExp {
        t.Errorf("Expected html=\n\"%v\"\nbut got\n\"%v\"", htmlStrExp, htmlStr)
    }
}

