#!/usr/bin/env bash


mode="${1:-report}"
report_file="report-$(date -u +"%Y%m%d%H%M%S").bin"

echo "Starting status checks (${mode})… Please wait."
/usr/bin/env vegeta attack -duration=60s -rate=5 -targets=targets.list -output="${report_file}"


case "${mode}" in
  "plot")
    /usr/bin/env vegeta plot -title="Status Checks" "${report_file}" > results.html
    /usr/bin/open results.html
    ;;
  "report")
    /usr/bin/env vegeta report "${report_file}"
    ;;
  *)
    echo "Usage: ${0} (plot|report|)"
    exit 1
esac

