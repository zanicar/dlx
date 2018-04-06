# dlx
A Golang implementation of Dancing Links (Algorithm X) as described by Donald Knuth.

Donald Knuth, 2000. Dancing Links. Millennial Perspectives in Computer Science, p159, Volume 187. The paper was accessed from: https://arxiv.org/pdf/cs/0011047v1.pdf

This was my first proper Golang endeavor to familiarize myself with the language and only some minor changes have been made over the years since it was originally written. The implementation tries to follow the paper as closely as possible according to the Golang specifications. It is a standalone package that should be able to solve any exact cover problem.

## Package Basics
The package contains three types: Matrix, Element and Head.

Matrix is a container of Element and these are modelled after the Golang standard library's List. Matrix has two separate classes of elements, namely header elements and standard elements. The header elements are also used to serve the purpose of sentinel elelents, thus Matrix does not contain explicit sentinel elements. Head is used as the Value of the elements contained in the Matrix's header. The header is used to define the various constraint columns within the sparse matrix.

The API allows for the creation of a newly initialized Matrix, for pushing Header Elements and Standard Elements with PushHead and PushItem respectively. Solve can be invoked to find all solutions for the given problem in its current state, it returns the found solutions in a slice, each being a slice of strings that exactly covers the problem space.

## Package Documentation
https://godoc.org/github.com/zanicar/dlx

## Sudoku Implementation
A terminal based Sudoku solver that supports multiple sizes and multiple solutions: github.com/zanicar/sudoku
