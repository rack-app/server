./list.sh | awk '{print($2)}' | xargs kill -9
find ./tmp -type s -delete