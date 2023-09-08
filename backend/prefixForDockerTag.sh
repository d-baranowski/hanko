export branch=$(git rev-parse --abbrev-ref HEAD)
export hash=$(git rev-parse --short HEAD)
export timestamp=$(date +%s)
git describe --exact-match --tags >/dev/null 2>&1 && git describe --exact-match --tags || echo $timestamp-$branch-$hash
