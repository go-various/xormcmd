#!/bin/bash
cmd=build/bin/xorm
function model() {
    $cmd reverse mysql "user:password@tcp(localhost)/test" templates/goxorm build/models $1
}
function view() {
    $cmd view mysql "user:password@tcp(localhost)/test" templates/goview build/views $1
}

case $1 in
  model)
  model "$2"
  ;;
  view)
    view "$2"
    ;;
  *)
    echo "Usage: $0 [model|view] [table_filters]"
esac
