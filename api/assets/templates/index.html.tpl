<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 3.2 Final//EN">
<html>
 <head>
  <title>Index of {{.Prefix}}</title>
 </head>
 <body>
<h1>Index of {{.Prefix}}</h1>
  <table>
   <tr><th valign="top"><img src="/assets/icons/blank.png" alt="[ICO]"></th><th><a href="?C=N;O=D">Name</a></th><th><a href="?C=M;O=A">Last modified</a></th><th><a href="?C=S;O=A">Size</a></th><th><a href="?C=D;O=A">User Meta</a></th></tr>
   <tr><th colspan="5"><hr></th></tr>
{{if .ParentPrefix}}
<tr><td valign="top"><img src="/assets/icons/back.png" alt="[PARENTDIR]"></td><td><a href="/index{{.ParentPrefix}}">Parent Directory</a></td><td>&nbsp;</td><td align="right">  - </td><td>&nbsp;</td></tr>
{{end}}
{{range .Files}}
	{{if .Dir}}
	<tr><td valign="top"><img src="/assets/icons/{{.Icon}}" alt="{{.IconAlt}}"></td><td><a href="/index{{.Path}}">{{.Name}}</a></td><td align="right">&nbsp;</td><td align="right">&nbsp;</td><td>&nbsp;</td></tr>
	{{else}}
	<tr><td valign="top"><img src="/assets/icons//{{.Icon}}" alt="{{.IconAlt}}"></td><td><a href="/api/v1/content{{.Path}}">{{.Name}}</a></td><td align="right">{{.LastModified}}</td><td align="right">{{.Size}}</td><td>{{.UserMeta}}</td></tr>
	{{end}}
{{end}}
   <tr><th colspan="5"><hr></th></tr>
</table>
</body>
</html>
