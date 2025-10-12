# Data Peringatan Dini Cuaca Terbuka BMKG (CAP)

Data Data Peringatan Dini Cuaca Terbuka BMKG telah tersedia di portal https://data.bmkg.go.id/peringatan-dini-cuaca dengan format XML Common Alerting Protocol (CAP). Berikut kode baris pemrograman PHP yang digunakan dalam mengolah data peringatan dini cuaca tersebut.

## Mengolah Data CAP Nowcast RSS Feed

#### Kode Baris PHP untuk Mengolah Data `nowcast-feed.php`
```php
<?php
// URL to the CAP Alert Nowcast RSS feed
$rss_url = "https://www.bmkg.go.id/alerts/nowcast/id/rss.xml"; // for Indonesian
// Or use https://www.bmkg.go.id/alerts/nowcast/en/rss.xml for English

// Get RSS content from the URL
$rss_content = @file_get_contents($rss_url);

// Check if fail to get content
if ($rss_content === false) {
    die("ERROR: Gagal mengambil data.");
}

// Load the XML string using SimpleXMLElement
libxml_use_internal_errors(true);
$xml = simplexml_load_string($rss_content);

// Check if XML loading failed
if ($xml === false) {
    $errors = libxml_get_errors();
    $error_messages = [];
    foreach ($errors as $error) {
        $error_messages[] = $error->message;
    }
    libxml_clear_errors();
    die(
        "ERROR: Gagal memuat file XML. " .
            htmlspecialchars(implode(", ", $error_messages))
    );
}

// Set header
header("Content-Type: text/html; charset=utf-8");
?>
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Peringatan Dini Cuaca BMKG</title>
    <style>
        body { font-family: sans-serif; line-height: 1.5; padding: 15px; }
        h2, h3, h4 { margin-top: 1.5em; margin-bottom: 0.5em; }
        ul { list-style: none; padding-left: 0; }
        li { margin-bottom: 1em; border-bottom: 1px solid #eee; padding-bottom: 1em; }
        .item-meta { font-size: 0.9em; color: #555; }
        pre { background-color: #f4f4f4; padding: 10px; border: 1px solid #ddd; overflow-x: auto; }
    </style>
</head>
<body>

<h1>Peringatan Dini Cuaca BMKG (CAP Alert Nowcast)</h1>

<?php if (isset($xml->channel->item)) {
    echo "<h2>Daftar Peringatan Aktif:</h2>";
    echo "<ul>";
    foreach ($xml->channel->item as $item) {
        $title = isset($item->title)
            ? htmlspecialchars((string) $item->title)
            : "Tanpa Judul";
        $link = isset($item->link)
            ? htmlspecialchars((string) $item->link)
            : "#";
        $description = isset($item->description)
            ? htmlspecialchars((string) $item->description)
            : "Tanpa Deskripsi";
        $author = isset($item->author)
            ? htmlspecialchars((string) $item->author)
            : "N/A";
        $pubDate = isset($item->pubDate)
            ? htmlspecialchars((string) $item->pubDate)
            : "N/A";

        echo "<li>";
        echo "<h3><a href=\"" .
            $link .
            "\" target=\"_blank\">" .
            $title .
            "</a></h3>";
        echo "<p class=\"item-meta\"><strong>Dipublikasikan pada:</strong> " .
            $pubDate .
            " | <strong>Diterbitkan oleh:</strong> " .
            $author .
            "</p>";
        echo "<p>" . nl2br($description) . "</p>";
        echo "</li>";
    }
    echo "</ul>";
} else {
    echo "<p>Tidak ada peringatan dini cuaca yang ditemukan dalam RSS Feed.</p>";
}
// Debugging $xml
/*
echo "<pre>";
print_r($xml);
echo "</pre>";
*/
?>
</body>
</html>
```
## Attribution / Sumber Data
**Perhatian!** Wajib untuk mencantumkan BMKG (Badan Meteorologi, Klimatologi, dan Geofisika) sebagai sumber data dan menampilkannya pada aplikasi atau sistem Anda.
