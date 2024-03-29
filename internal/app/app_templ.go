// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package app

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func hello() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
			if !templ_7745c5c3_IsBuffer {
				templ_7745c5c3_Buffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<h2 class=\"text-center text-neutral-800 dark:text-neutral-200\">Try again Sucka!!</h2>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = Base().Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func Base() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var3 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var3 == nil {
			templ_7745c5c3_Var3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<html><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><link rel=\"stylesheet\" href=\"https://cdn.jsdelivr.net/gh/iconoir-icons/iconoir@main/css/iconoir.css\"><script src=\"https://cdn.tailwindcss.com\"></script><script>\n        if (localStorage.theme === 'dark' || (!('theme' in localStorage) \n            && window.matchMedia('(prefers-color-scheme: dark)').matches)) {\n            document.documentElement.classList.add('dark');\n        } else {\n            document.documentElement.classList.remove('dark');\n        }\n        // whenever the user explicitly chooses light mode\n        localStorage.theme = 'light'\n        // whenever the user chooses dark mode\n        localStorage.theme = 'dark'\n        // OS preference\n        localStorage.removeItem('theme')\n    </script><script src=\"https://unpkg.com/htmx.org@1.9.10\" integrity=\"sha384-D1Kt99CQMDuVetoL1lrYwg5t+9QdHe7NLX/SoJYkXDFfX37iInKRy5xLSi8nO7UC\" crossorigin=\"anonymous\"></script><script src=\"https://unpkg.com/htmx.org/dist/ext/json-enc.js\"></script><style>\n    .htmx-indicator{\n        display:none;\n    }\n    .htmx-request .htmx-indicator{\n        display:inline;\n    }\n    .htmx-request.htmx-indicator{\n        display:inline;\n    }\n    </style></head><body class=\"h-full bg-neutral-100 dark:bg-neutral-800 text-neutral-800 dark:text-neutral-200 px-2\"><div class=\"flex flex-col\"><h1 class=\"text-center text-3xl my-5\">LED Lights</h1><p class=\"text-center text-md text-neutral-300\" id=\"mode\">off</p></div><div class=\"container mx-auto px-3 py-4 w-half bg-neutral-200 dark:bg-neutral-700 rounded-sm flex flex-col\"><div class=\"w-half\"><span class=\"isolate inline-flex rounded-2xl mx-auto\"><button type=\"button\" class=\"relative inline-flex items-center rounded-l-2xl bg-lime-100 px-3 py-2 text-sm font-semibold text-lime-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-10\" hx-get=\"/on\" hx-swap=\"none\">On</button> <button type=\"button\" class=\"relative inline-flex items-center bg-neutral-200 px-3 py-e text-sm font-semibold text-neutral-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-10\" hx-post=\"/reset\" hx-swap=\"none\">Reset</button> <button type=\"button\" class=\"relative -ml-px inline-flex items-center rounded-r-2xl bg-red-100 px-3 py-2 text-sm font-semibold text-red-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-10\" hx-get=\"/off\" hx-swap=\"none\">Off</button></span></div><div class=\"container mx-auto px-2 py-1 my-1 w-half\"><form hx-post=\"/mode\" hx-ext=\"json-enc\" hx-swap=\"innerHTML\" hx-target=\"#mode\" class=\"rounded-sm bg-neutral-300 px-5 py-2 mx-auto w-50 flex flex-col\"><div class=\"my-2\"><label for=\"mode\" class=\"block text-sm font-medium leading-6 text-gray-900\">Mode</label> <select id=\"mode\" name=\"mode\" class=\"mt-2 block w-full rounded-l border-0 py-1.5 pl-3 pr-10 text-gray-900 ring-1 ring-inset ring-gray-300 focus:ring-2 focus:ring-indigo-600 sm:text-sm sm:leading-6\"><option value=\"warm\">Warm</option> <option value=\"static\" selected>Static</option> <option value=\"rgbwave\">RGB Wave</option> <option value=\"rgbfade\">RGB Fade</option></select></div><input class=\"rounded-l my-2 px-2 py-1 text-neutral-800\" type=\"color\" name=\"rgbColor\" placeholder=\"0,0,0\"> <input class=\"rounded-l my-2 px-2 py-1 text-neutral-800\" type=\"brightness\" name=\"brightness\" placeholder=\"128\"> <button class=\"rounded-2xl bg-indigo-500 text-sm text-indigo-50 px-3 py-2 font-semibold my-1 items-center \n               ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-10\" type=\"submit\">Switch</button></form></div></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
