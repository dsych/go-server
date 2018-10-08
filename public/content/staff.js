window.addEventListener("load", () => {
    document.querySelector("button#search").addEventListener("click", () => {
        const firstName = document.querySelector("#firstName").value;
        const lastName = document.querySelector("#lastName").value;
        const employmentId = document.querySelector("#employmentId").value;
        const manager = document.querySelector("#manager").value;

        const baseUrl = `${window.location.protocol}//${window.location.host}`;
        fetch(`${baseUrl}/api/searchStaff`, {
            method: "POST",
            credentials: "include",
            body: JSON.stringify({
                firstName,
                lastName,
                employmentId,
                manager
            })
        })
            .then(res => {
                if (res.ok) {
                    return res.json();
                } else {
                    throw res.statusText;
                }
            })
            .then(json => {
                tableCreate(json);
            })
            .catch(err => {
                document.querySelector("#result").innerHTML = err;
                console.log(err);
            });
    });
});

const tableCreate = json => {
    const createCell = (value, td) => {
        td.appendChild(document.createTextNode(value));
        td.style.border = "1px solid black";
    };

    const tbl = document.createElement("table");
    tbl.id = "result";
    tbl.style.width = "100px";
    tbl.style.border = "1px solid black";
    const header = tbl.createTHead().insertRow();
    createCell("First Name", header.insertCell());
    createCell("Last Name", header.insertCell());
    createCell("Manager", header.insertCell());
    createCell("Gender", header.insertCell());
    createCell("Date of Birth", header.insertCell());
    createCell("Health Card Number", header.insertCell());
    createCell("SIN", header.insertCell());
    createCell("University", header.insertCell());
    createCell("Home Address", header.insertCell());
    createCell("Email", header.insertCell());
    createCell("Employment ID", header.insertCell());
    createCell("Job Role", header.insertCell());
    createCell("Pay", header.insertCell());

    for (let i = 0; i < json.length; i++) {
        const tr = tbl.insertRow();
        createCell(json[i]["firstName"], tr.insertCell());
        createCell(json[i]["lastName"], tr.insertCell());
        createCell(json[i]["manager"], tr.insertCell());
        createCell(json[i]["gender"], tr.insertCell());
        createCell(json[i]["DOB"], tr.insertCell());
        createCell(json[i]["healthCard"], tr.insertCell());
        createCell(json[i]["SIN"], tr.insertCell());
        createCell(json[i]["university"], tr.insertCell());
        createCell(json[i]["homeAddress"], tr.insertCell());
        createCell(json[i]["email"], tr.insertCell());
        createCell(json[i]["employmentId"], tr.insertCell());
        createCell(json[i]["jobRole"], tr.insertCell());
        createCell(json[i]["pay"], tr.insertCell());
    }

    // build a reference to the existing node to be replaced
    const sp2 = document.getElementById("result");
    const parentDiv = sp2.parentNode;

    // replace existing node sp2 with the new span element sp1
    parentDiv.replaceChild(tbl, sp2);
};
