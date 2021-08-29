function refreshRender() {
    let editId = document.getElementById("editid").innerText;
    console.log("Refreshing type for " + editId);
    renderEdit(editId);
}

function classifyEdit(classification, confirmation) {
    let editId = document.getElementById("editid").innerText;
    console.log("Classifying " + editId + " as " + classification + " (" + confirmation + ")");

    let req = new XMLHttpRequest();
    req.onreadystatechange = function(){
        if (this.readyState !== 4) {
            return;
        }

        if (this.status !== 200) {
            alert('Failed to classify edit');
            return;
        }

        // API wants to confirm - ask the user
        let require_confirmation = JSON.parse(this.responseText)["require_confirmation"];
        if (require_confirmation) {
            let confirmation = confirm("Are you sure?");
            if (!confirmation) {
                return;
            }
            return classifyEdit(classification, true);
        }

        // We are done - onto the next
        loadNextEditId();
    }
    req.open("POST", "/api/user-classification", true);
    req.send(JSON.stringify({
        "edit_id": parseInt(editId),
        "comment": document.getElementById("comment").value,
        "classification": classification,
        "confirmation": confirmation,
    }));
}

function renderEdit(editId) {
    let urlType = "n";
    document.getElementsByName("url_type").forEach(function(radio) {
        if (radio.checked) {
            urlType = radio.value;
        }
    });
    console.log("Rendering: " + editId + " using " + urlType);

    let url = "https://en.wikipedia.org/w/index.php?action=view&diff=" + editId;
    if (urlType === "d") {
        url = "https://en.wikipedia.org/w/index.php?action=view&diffonly=1&diff=" + editId;
    } else if (urlType === "r") {
        url = "https://en.wikipedia.org/w/index.php?action=render&diffonly=1&diff=" + editId;
    }
    document.getElementById("editid").innerText = editId;
    document.getElementById("iframe").setAttribute("src", url);
}

function loadNextEditId() {
    document.getElementById("comment").value = "";

    let req = new XMLHttpRequest();
    req.onreadystatechange = function(){
        if (this.readyState !== 4) {
            return;
        }

        if (this.status !== 200) {
            alert('Failed to retrieve pending edit');
            return;
        }

        let editId = JSON.parse(this.responseText)["edit_id"];
        renderEdit(editId);
    }
    req.open("GET", "/api/edit/next", true);
    req.send();
}

function loadDetails() {
    let editId = document.getElementById("editid").innerText;
    window.open("/admin/details/" + editId, true);
}

window.onload = function() {
    loadNextEditId();
}
