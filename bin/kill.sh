./list.sh | awk '{print($2)}' | xargs kill -9
[ -d ./tmp ] && find ./tmp -type s -delete