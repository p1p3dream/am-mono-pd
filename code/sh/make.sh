# This file MUST be POSIX-compliant for max compatibility with other shells.

ENV_FILE=""

while getopts "e:" opt; do
    case $opt in
        e)
            ENV_FILE="$OPTARG"
            ;;
        *)
            ;;
    esac
done

if [ -z "${ENV_FILE}" ]; then
    echo "error: -e is required."
    exit 1
fi

shift $((OPTIND - 1))

while IFS='=' read -r key value; do
    if [ -n "$key" ] && [ -n "$value" ]; then
        export "$key=$value"
    fi
done < "${ENV_FILE}"
