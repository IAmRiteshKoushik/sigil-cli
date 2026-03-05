build:
  @go build -ldflags="-X 'github.com/IAmRiteshKoushik/sigil/cmd.Version=$(git describe --tags)'" -o sigil main.go

# This requires you to install - github.com/pdfcpu/pdfcpu
# Works for both participation and recognition certificates
stamp-participation name event input output nameoffset="-115" eventnameoffset="-35":
  @echo "Step 1: Stamping Name '{{name}}'..."
  @pdfcpu stamp add -mode text -- \
      "{{name}}" \
      "scale:2.5 abs, pos:bc, off:{{nameoffset}} 1600, al:c, rot:0, fillc:#000000, font:Helvetica-Bold" \
      "{{input}}" "temp_intermediate.pdf"

  @echo "Step 2: Stamping event '{{event}}'"
  @pdfcpu stamp add -mode text -- \
      "{{event}}" \
      "scale:2.5 abs, pos:bc, off:{{eventnameoffset}} 1460, rot:0, fillc:#000000, font:Helvetica-Bold" \
      "temp_intermediate.pdf" "{{output}}"

  @rm temp_intermediate.pdf
  @echo "Success! Created {{output}}"

# Works for winner certificates
stamp-winner name event pos input output:
  @echo "Step 1: Stamping Name '{{name}}'..."
  @pdfcpu stamp add -mode text -- \
      "{{name}}" \
      "scale:3.5 abs, pos:bl, off:1250 1570, rot:0, fillc:#000000, font:Helvetica" \
      "{{input}}" "temp_intermediate.pdf"

  @echo "Step 2: Stamping event '{{event}}'"
  @pdfcpu stamp add -mode text -- \
      "{{event}}" \
      "scale:3.5 abs, pos:bl, off:1250 1420, rot:0, fillc:#000000, font:Helvetica" \
      "temp_intermediate.pdf" "{{output}}"

  @rm temp_intermediate.pdf
  @echo "Success! Created {{output}}"
