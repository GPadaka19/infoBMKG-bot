<?php
// URL to the CAP Alert XML
$cap_url = "https://www.bmkg.go.id/alerts/nowcast/id/CJB20251013001_alert.xml"; // for Indonesian
// Or use https://www.bmkg.go.id/alerts/nowcast/id/CJB20251013001_alert.xml for English

// Get CAP Alert content from the URL
$xml_string = @file_get_contents($cap_url);

// Check if fail to get content
if ($xml_string === false) {
    die("ERROR: Gagal mengambil data.");
}

// Load the XML string using SimpleXMLElement
libxml_use_internal_errors(true);
$xml = simplexml_load_string($xml_string);

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

// Register CAP namespace
$namespaces = $xml->getNamespaces(true);
$cap_ns = isset($namespaces[""])
    ? $namespaces[""]
    : "urn:oasis:names:tc:emergency:cap:1.2";

// Set header
header("Content-Type: text/html; charset=utf-8");

// Accessing <info> block
$info = $xml->children($cap_ns)->info;
$polygons_data = [];

if ($info) {
    // Extract polygon data
    if (isset($info->area)) {
        foreach ($info->area as $area) {
            if (isset($area->polygon)) {
                foreach ($area->polygon as $polygon_str) {
                    $polygon_coords = [];
                    // Split the string into coordinate pairs
                    $coord_pairs = explode(" ", (string) $polygon_str);
                    foreach ($coord_pairs as $pair) {
                        // Split each pair into lat and lon
                        $lat_lon = explode(",", $pair);
                        if (
                            count($lat_lon) === 2 &&
                            is_numeric($lat_lon[0]) &&
                            is_numeric($lat_lon[1])
                        ) {
                            // Leaflet expects [latitude, longitude]
                            $polygon_coords[] = [
                                (float) $lat_lon[0],
                                (float) $lat_lon[1],
                            ];
                        }
                    }
                    if (!empty($polygon_coords)) {
                        $polygons_data[] = $polygon_coords;
                    }
                }
            }
        }
    }
}
?>
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Detail Peringatan Dini Cuaca</title>
    <!-- LeafletJS CSS -->
    <link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" integrity="sha256-p4NxAoJBhIIN+hmNHrzRCf9tD/miZyoHS5obTRR9BMY=" crossorigin=""/>
    <!-- LeafletJS JS -->
    <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js" integrity="sha256-20nQCchB9co0qIjJZRGuk2/Z9VM+kNiyxNV1lvTlZBo=" crossorigin=""></script>
    <style>
        body { font-family: sans-serif; line-height: 1.5; padding: 15px; }
        p { margin: 0.5em 0; }
        strong { display: inline-block; width: 120px; }
        img.image { max-width: 100%; height: auto; border: 1px solid #ddd; padding: 5px; background-color: #f9f9f9; }
        #map { height: 400px; width: 100%; margin-top: 20px; border: 1px solid #ccc; }
    </style>
</head>
<body>

<?php if ($info) {
    $headline = isset($info->headline)
        ? htmlspecialchars((string) $info->headline)
        : "Tanpa Judul Headline";
    $event = isset($info->event)
        ? htmlspecialchars((string) $info->event)
        : "N/A";
    $urgency = isset($info->urgency)
        ? htmlspecialchars((string) $info->urgency)
        : "N/A";
    $severity = isset($info->severity)
        ? htmlspecialchars((string) $info->severity)
        : "N/A";
    $certainty = isset($info->certainty)
        ? htmlspecialchars((string) $info->certainty)
        : "N/A";
    $effective = isset($info->effective)
        ? htmlspecialchars((string) $info->effective)
        : "N/A";
    $expires = isset($info->expires)
        ? htmlspecialchars((string) $info->expires)
        : "N/A";
    $senderName = isset($info->senderName)
        ? htmlspecialchars((string) $info->senderName)
        : "N/A";
    $description = isset($info->description)
        ? htmlspecialchars((string) $info->description)
        : "N/A";
    $web = isset($info->web) ? htmlspecialchars((string) $info->web) : "#";
    $areaDesc = isset($info->area->areaDesc)
        ? htmlspecialchars((string) $info->area->areaDesc)
        : "N/A";

    echo "<h1>" . $headline . "</h1>";

    echo "<p><strong>Jenis Peringatan:</strong> " . $event . "</p>";
    echo "<p><strong>Urgency:</strong> " . $urgency . "</p>";
    echo "<p><strong>Severity:</strong> " . $severity . "</p>";
    echo "<p><strong>Certainty:</strong> " . $certainty . "</p>";
    echo "<p><strong>Dimulai pada:</strong> " . $effective . "</p>";
    echo "<p><strong>Berakhir pada:</strong> " . $expires . "</p>";
    echo "<p><strong>Dipublikasikan oleh:</strong> " . $senderName . "</p>";
    echo "<p><strong>Deskripsi:</strong> " . nl2br($description) . "</p>";
    echo "<p><strong>Area Terdampak:</strong> " . $areaDesc . "</p>";

    // Map Container
    echo '<div id="map"></div>';

    // Infografik
    echo '<p><strong>Infografik:</strong> <img src="' .
        $web .
        '" alt="Infografik" class="image"></p>';
} else {
    echo "<p>Blok 'info' tidak ditemukan dalam CAP Alert.</p>";
}
// Debugging $xml
/*
echo "<pre>";
print_r($xml);
echo "</pre>";
*/
?>

<script>
document.addEventListener('DOMContentLoaded', function () {
    // Convert PHP polygon data to JavaScript
    const polygonsData = <?= json_encode($polygons_data) ?>;

    if (polygonsData.length > 0) {
        // Initialize the map. Use the first point of the first polygon as a rough center.
        const map = L.map('map').setView(polygonsData[0][0], 8);

        // Add a tile layer from OpenStreetMap
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        }).addTo(map);

        // Create a feature group to hold all polygons
        const polygonGroup = L.featureGroup();

        // Add each polygon to the map
        polygonsData.forEach(polygonCoords => {
            const polygon = L.polygon(polygonCoords, {
                color: 'red',
                fillColor: 'orange',
                fillOpacity: 0.5,
                weight: 1
            });
            polygonGroup.addLayer(polygon);
        });

        // Add the group to the map
        polygonGroup.addTo(map);

        // Fit the map to the bounds of all polygons
        map.fitBounds(polygonGroup.getBounds());
    }
});
</script>

</body>
</html>
