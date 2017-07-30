const
    express = require("express"),
    // https = require("https"),
    path = require("path");

    var app = express();

    // Main Page
    app.use("/", express.static(path.join(__dirname, "public")));
    
    // 404
    app.use(function(req, res) {
        res.status(404);
        res.sendFile(path.join(__dirname, "public/404.html"));
    })

    app.listen(80, () => console.log("Started server on port 80"));