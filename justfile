build:
  @go build -ldflags="-X 'github.com/IAmRiteshKoushik/sigil/cmd.Version=$(git describe --tags)'" -o sigil main.go

# This requires you to install - github.com/pdfcpu/pdfcpu

# Setup the nameoffset, eventnameoffset and posoffset based on where you want to 
# position the text on the canvas of the PDF. This requires a lot of trial and 
# error before you get it right. Don't change other configuration parameters as 
# it is already in the best config for managing things. x = 0 is the center line
# passing through a certificate. Set your x-offset based on that. For your 
# y-offset, you are on your own :)

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
stamp-winner name event pos input output nameoffset="0" eventnameoffset="0" posoffset="-200":
  @echo "Step 1: Stamping Name '{{name}}'..."
  @pdfcpu stamp add -mode text -- \
      "{{name}}" \
      "scale:2.5 abs, pos:bc, off:{{nameoffset}} 1600, al:c, rot:0, fillc:#000000, font:Helvetica-Bold" \
      "{{input}}" "step_1.pdf"

  @echo "Step 2: Stamping event '{{event}}'"
  @pdfcpu stamp add -mode text -- \
      "{{event}}" \
      "scale:2.5 abs, pos:bc, off:{{eventnameoffset}} 1460, rot:0, fillc:#000000, font:Helvetica-Bold" \
      "step_1.pdf" "step_2.pdf"

  @echo "Step 3: Stamping position '{{pos}}'"
  @pdfcpu stamp add -mode text -- \
      "{{pos}}" \
      "scale:2.5 abs, pos:bc, off:{{posoffset}} 1260, rot:0, fillc:#000000, font:Helvetica-Bold" \
      "step_2.pdf" "{{output}}"

  @rm step_1.pdf
  @rm step_2.pdf
  @echo "Success! Created {{output}}
