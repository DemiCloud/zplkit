^XA
^FX asset-tag/2.25x1.25/001 - 2.25 x 1.25 in @ 203 dpi: company name, human-readable tag, Code 128, QR to asset URL
^CI28
^PW456
^LL254
^LH0,0

^FX --- Branding ---
^FO38,40
^A0N,30,30
^FD{{.CompanyName}} Asset^FS

^FX --- Human-readable asset tag ---
^FO38,80
^A0N,45,38
^FD{{.Tag}}^FS

^FX --- Code 128 asset tag ---
^FO38,141
^BY2,2,72
^BCN,72,N,N,N,A
^FD{{.Tag}}^FS

^FX --- QR direct asset URL ---
^FO305,44
^BQN,2,4
^FDQA,{{.AssetURL}}{{.Tag}}^FS

^XZ
