#!/usr/bin/zsh -e

setopt pipefail

PREFIX="/${ABODEMINE_NAMESPACE}/"
VERSION="$1"

case "${VERSION}" in
   v2)
        aws ssm \
            get-parameters-by-path \
            --path "${PREFIX}" \
            --recursive \
            --query "Parameters[*].{Name:Name,Value:Value}" \
        | jq '
            reduce .[] as $item (
                {};
                if $item.Name | endswith(".JSON") then
                    .[$item.Name | rtrimstr(".JSON") | sub("^'${PREFIX}'"; "")] = ($item.Value | fromjson)
                else
                    .[$item.Name | sub("^'${PREFIX}'"; "")] = $item.Value
                end
            )
        '
        ;;
    *)
        aws ssm \
            get-parameters-by-path \
            --path "${PREFIX}" \
            --recursive \
            --query "Parameters[*].{Name:Name,Value:Value}" \
        | jq '
            reduce .[] as $item (
                {};
                setpath(
                    (
                        $item.Name
                        | sub("^'${PREFIX}'"; "")
                        | split("/")
                    );
                    $item.Value
                )
            )
        '
        exit 0
        ;;
esac
