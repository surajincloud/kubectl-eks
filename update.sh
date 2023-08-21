#!/bin/bash

version="$1"

url="https://github.com/surajincloud/kubectl-eks/releases/download/v${version}/kubectl-eks_${version}_checksums.txt"

echo "downloaded file"
echo $url

# Download the file
wget -q "$url" -O checksums.txt

output="plugins/kubectl-eks.yaml"

template="template.yaml"

sed "s/{{version}}/$version/g" "$template" > "$output"

while IFS=" " read -r sha256 filename; do
  os_arch=$(echo "$filename" | cut -d "_" -f 3,4 | cut -d "." -f 1)
  os=$(echo "$os_arch" | cut -d "_" -f 1)
  arch=$(echo "$os_arch" | cut -d "_" -f 2)
  bin=$(echo "$filename" | cut -d "_" -f 1)

  yaml_entry="\
  - selector:
      matchLabels:
        os: $os
        arch: $arch
    uri: https://github.com/surajincloud/kubectl-eks/releases/download/v${version}/${filename}
    sha256: $sha256
    bin: $bin"

  echo "$yaml_entry" >> "$output"
done < checksums.txt

rm checksums.txt

echo "Generated YAML content in $output"
