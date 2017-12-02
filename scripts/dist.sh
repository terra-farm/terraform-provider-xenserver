#/bin/sh

mkdir dist
OS="linux"
case "${OSTYPE}" in
  darwin*)  OS="darwin" ;; 
  linux*)   OS="linux" ;;
  *)        echo "unknown: ${OSTYPE}" ;;
esac
ARCH=$(uname -m)
cp "${GOPATH}/bin/terraform-provider-xenserver" "dist/terraform-provider-xenserver-${TRAVIS_TAG}-${OS}-${ARCH}"
md5sum "dist/terraform-provider-xenserver-${TRAVIS_TAG}-${OS}-${ARCH}" | awk '{ print $1 }' > "dist/terraform-provider-xenserver-${TRAVIS_TAG}-${OS}-${ARCH}.md5sum"
ls -lsa dist/
cat "dist/terraform-provider-xenserver-${TRAVIS_TAG}-${OS}-${ARCH}.md5sum"
