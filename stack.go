package main

import (
	"errors"
)

type Stack[T any] struct {
	elements []T
}

func (s *Stack[T]) Push(value T) {
	s.elements = append(s.elements, value)
}

func (s *Stack[T]) Pop() (T, error) {
	var zero T
	if s.IsEmpty() {
		return zero, errors.New("stack is empty")
	}

	// Get index of the top element
	index := len(s.elements) - 1
	value := s.elements[index]

	// Prevent memory leaks by zeroing out the reference before re-slicing
	s.elements[index] = zero
	s.elements = s.elements[:index]

	return value, nil
}

func (s *Stack[T]) Peek() (T, error) {
	var zero T
	if s.IsEmpty() {
		return zero, errors.New("stack is empty")
	}
	return s.elements[len(s.elements)-1], nil
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.elements) == 0
}

func (s *Stack[T]) Size() int {
	return len(s.elements)
}
