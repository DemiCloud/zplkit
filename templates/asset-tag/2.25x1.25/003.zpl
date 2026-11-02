^XA
^FX asset-tag/2.25x1.25/003 - 2.25 x 1.25 in @ 203 dpi: rounded border, company name, human-readable tag, Code 128, QR to asset URL
^CI28
^PW456
^LL254
^LH0,0

^FX --- Border with rounded corners ---
^FO12,12
^GB432,230,4,B,3^FS

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

^FX --- QR direct asset URL --- (^BY2,2,0 resets the barcode's height,
^FX which ^BQ otherwise inherits as an unwanted Y-offset) ---
^BY2,2,0
^FO306,76
^BQN,2,3
^FDQA,{{.AssetURL}}{{.Tag}}^FS

^XZ
