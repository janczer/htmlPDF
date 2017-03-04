## Package for create pdf form html

   Implement series article [Let's build a toy browser engine!](https://limpet.net/mbrubeck/2014/08/08/toy-layout-engine-1.html)

Install:

    go get github.com/janczer/htmlPDF
    go get github.com/jung-kurt/gofpdf

Exemple:
First you need create 3 files. With html tags, css style and Go code:

`first.html`:

    <div class="a">
        <div class="b">
            <div class="c">
                <div class="d">
                    <div class="e">
                        <div class="f">
                            <div class="g">
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
	</div>

`style.css`:

    * { display: block; padding: 5px; }
    .a { background: #ff0000; }
    .b { background: #ffa500; }
    .c { background: #ffff00; }
    .d { background: #008000; }
    .e { background: #0000ff; }
    .f { background: #4b0082; }
    .g { background: #800080; }
    
    
`main.go`:

    package main

    import "github.com/janczer/htmlPDF"

    func main() {
        htmlPDF.Generate("first.html", "style.css", "hello.pdf")
    }


### Todo:
<hr>

- [ ] Add support Anonymous block
- [ ] Add support Inline block
- [ ] Add use CSS from style tag `<style></style>`

<hr>
