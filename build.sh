#!/bin/bash
cmd=build/bin/xorm
function model() {
    $cmd reverse mysql "root:123456@tcp(localhost)/hytest" templates/goxorm build/models $1
}
function view() {
    $cmd view mysql "root:123456@tcp(localhost)/hytest" templates/goview build/views $1
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
