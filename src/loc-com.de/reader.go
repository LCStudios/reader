package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"loc-com.de/feed"
	"log"
	"net/http"
	"os"
)

// /* Print something */
// func main() {
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		page, _ := ioutil.ReadFile("index.html")
// 		fmt.Fprint(w, string(page))
// 	})
// 	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
// 	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
// 	log.Fatal(http.ListenAndServe(":8000", nil))
// }

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadFile("index.html")
		fmt.Fprint(w, string(body))
	})
	http.HandleFunc("/rss", getRSSFeed)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func getRSSFeed(w http.ResponseWriter, r *http.Request) {

	response, err := http.Get("http://loc-blog.de/rss.php?blog_id=5")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	feed, _ := feed.Decode([]byte(contents))
	respJSON, _ := json.Marshal(feed)
	fmt.Fprint(w, string(respJSON))
}

//  contents := `<?xml version="1.0" encoding="UTF-8"?>
// <rss version="2.0"
//     xmlns:content="http://purl.org/rss/1.0/modules/content/"
//     xmlns:wfw="http://wellformedweb.org/CommentAPI/"
//     xmlns:dc="http://purl.org/dc/elements/1.1/"
//     xmlns:atom="http://www.w3.org/2005/Atom"
//     xmlns:sy="http://purl.org/rss/1.0/modules/syndication/"
//     xmlns:slash="http://purl.org/rss/1.0/modules/slash/"
//     >    <channel>
//         <title>LocCom Blog</title>
//         <description>LocCom, Community, Dev-Updates</description>
//         <link>http://www.loc-blog.de/index.php?site=blog&amp;blog_id=5</link>
//         <lastBuildDate>Sat, 13 Oct 2012 01:53:23 +0200</lastBuildDate>
//         <generator>FeedCreator 1.7.2-ppt (info@mypapit.net)</generator>
//         <item>
//             <title>Neuen Blog erstellen - HowTo</title>
//             <link>http://www.loc-blog.de/index.php?site=blog_post&amp;post_id=50</link>
//             <description></description>
//             <enclosure url="http://www.scripting.com/mp3s/weatherReportSuite.mp3" length="12216320" type="audio/mpeg" />
//             <content:encoded><![CDATA[<p>Hier eine kurze Erkl&auml;rung zum erstellen von Blogs.</p>
// <p>Als erstes klickt ihr auf "+ Neuen Blog erstellen".</p>
// <p><img style="vertical-align: middle;" title="screen-capture-3.jpg" src="files/471349828327.jpg" alt="screen-capture-3.jpg" width="397" height="117" /></p>
// <p>Danach kommt ihr auf folgende Seite:</p>
// <p><img title="screen-capture-4.jpg" src="files/481349828327.jpg" alt="screen-capture-4.jpg" width="607" height="200" /></p>
// <p>Bei "Blog Titel" tragt ihr einfach den Namen des Blogs ein. Dieser Titel wird dann im Browser als Title angezeigt und in deinem ganzen Blog als &Uuml;berschrift.<br />Die &Uuml;berschrift sollte den Inhalt des Blogs grob beschreiben.<br />Die Unter&uuml;berschrift dient dazu, die &Uuml;berschrift m&ouml;glicherweise zu erg&auml;nzen. Oftmals reicht der Platz des Titels nicht um ihn detailliert zu beschreiben. Dann kann man ihn so erweitern.</p>]]></content:encoded>
//             <pubDate>Wed, 10 Oct 2012 00:27:32 +0200</pubDate>
//         </item>
//         <item>
//             <title>Bilder Hochladen und Posten - HowTo</title>
//             <link>http://www.loc-blog.de/index.php?site=blog_post&amp;post_id=39</link>
//             <description></description>
//             <content:encoded><![CDATA[<p>Hallo,<br />ich m&ouml;chte nur ganz kurz erkl&auml;ren wie das posten von Bildern funktioniert.</p>
// <p>Fangen wir im Dashboard an:<br /><img title="Dashboard" src="files/411349208890.jpg" alt="Dashboard" width="612" height="116" /><br />Dort gibt es den Link "Dateiliste", &uuml;ber welchen du in eine &Uuml;bersicht aller Fotos kommst, die du f&uuml;r diesen Blog hochgeladen hast.<br />Das ganze sieht so aus:</p>
// <p><img title="dateiliste.jpg" src="files/421349208890.jpg" alt="dateiliste.jpg" width="604" height="199" /><br />Ganz oben hast du die M&ouml;glichkeit Fotos auf den Server hochzuladen. Einfach ausw&auml;hlen und hochladen. Sie erscheinen dann unten in einer Tabelle mit kleinem Vorschaubild, Namen und einem weiteren Link.<br />Das ist die Zentrale f&uuml;r alle Dateien, die du f&uuml;r den jeweiligen Blog hochgeladen hast.<br />&Uuml;ber den Link "Foto Posten" kannst du direkt ein Foto posten, wodurch der Editor direkt mit eingef&uuml;gtem Bild ge&ouml;ffnet wird.<br /><img title="new_post2.png" src="files/441349213753.png" alt="new_post2.png" width="600" height="226" /><br />Dort kannst du dann, durch Klicken auf das Bild, die Gr&ouml;&szlig;e ver&auml;ndern und es optimal anpassen.<br />Wenn du das Bild anklickst und dann auf das kleine B&auml;umchen (oben bei den Icons des Editors) klickst, bekommst du noch mehr M&ouml;glichkeiten dein Bild anzupassen.<br />Das Bild sollte generell nicht breiter sein als das Eingabefenster, da es sonst im sp&auml;teren Blog nicht besonders gut aussieht!</p>
// <p>Wenn man noch zus&auml;tzliche Bilder einf&uuml;gen m&ouml;chte, klickt man auf die Stelle wo das Foto hin soll und dann nochmals auf das B&auml;umchen. In dem Fenster das dann erscheint gibt es eine "Image List". Dort sind alle deine Fotos aufgef&uuml;hrt die du f&uuml;r diesen Blog hochgeladen hast. W&auml;hle einfach eins aus, passe es an und gehe dann auf "Insert".</p>
// <p>So, das wars soweit.<br />Wenn es noch Fragen gibt, schreibe sie einfach als Kommentar unten drunter :-)</p>]]></content:encoded>
//             <pubDate>Tue, 02 Oct 2012 20:33:37 +0200</pubDate>
//         </item>
//         <item>
//             <title>News auf loc-com.de überarbeitet</title>
//             <link>http://www.loc-blog.de/index.php?site=blog_post&amp;post_id=36</link>
//             <description></description>
//             <content:encoded><![CDATA[<p>Ein paar hundert Zeilen Code sp&auml;ter sind jetzt auch die News komplett objektorientiert geschrieben und keine vereinzelten mysql-Statements mehr &uuml;brig.</p>
// <p>Seit einer Weile wird LocCom in gro&szlig;en Teilen im Hintergrund umgebaut, da sich durch das ver&ouml;ffentlichen vieler neuer Features in kurzer Zeit und ohne gro&szlig;e Planung ein gro&szlig;er Haufen Altlasten angesammelt hatte und viel mit Copy &amp; Paste geschrieben wurde. Deswegen haben wir nun fast das ganze Backend neu geschrieben, in einem aufger&auml;umten objektorientierten Stil.</p>
// <p>Wenn das alles abgeschlossen ist k&ouml;nnen wir verschiedene neue Projekte und Features angehen. Ein Parallelprojekt ist ja bereits mit dem Blog hier in den Startl&ouml;chern. Auch gibt es bereits Prototypen f&uuml;r ein Dateisystem, und Events werden auch bald in LocCom Einzug halten.</p>
// <p>In der Zwischenzeit m&uuml;ssen allerdings noch einige Bugs gefixt werden und auch noch ein paar Stellen aufger&auml;umt werden.</p>
// <p>Freut euch auf die Dinge, die sehr bald kommen, aber wir k&ouml;nnen euch jetzt schon versichern , dass wir noch viele weitere Ideen haben, die wir euch hoffentlich auch in nicht allzu langer Zeit pr&auml;sentieren k&ouml;nnen.</p>
// <p>&nbsp;</p>
// <p>Euer LocCom-Team</p>]]></content:encoded>
//             <pubDate>Sun, 30 Sep 2012 22:48:42 +0200</pubDate>
//         </item>
//     </channel>
// </rss>
// `
