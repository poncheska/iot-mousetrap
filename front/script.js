let socket = new WebSocket("wss://https://smart-mousetrap.herokuapp.com");
function buttonDefinition(modal, modalOverlay, closeButton, openButton) {
	closeButton.addEventListener("click", function() {
		modal.classList.toggle("closed");
		modalOverlay.classList.toggle("closed");
        modal.querySelector("#Email").value = "";
        modal.querySelector("#Password").value = "";
	});

    openButton.addEventListener("click", function() {
        modal.classList.toggle("closed");
        modalOverlay.classList.toggle("closed");
    });
    }
    let modalSignIn = document.querySelector("#modal-sign-in"),
        modalJoin = document.querySelector("#modal-join"),
        modalOverlay = document.querySelector("#modal-overlay"),
        closeButtonSignIn = document.querySelector("#close-button-sign-in"),
        closeButtonJoin = document.querySelector("#close-button-join"),
        openButtonSignIn = document.querySelector("#open-button-sign-in");
        openButtonJoin = document.querySelector("#open-button-join");
    buttonDefinition(modalSignIn, modalOverlay, closeButtonSignIn, openButtonSignIn);
    buttonDefinition(modalJoin, modalOverlay, closeButtonJoin, openButtonJoin);

async function submit() {
	let signInFormData = new FormData(document.querySelector("#sign-in-form"));
	let joinFormData = new FormData(document.querySelector("#join-form"));
	modalSignIn.classList.toggle("closed");
	modalOverlay.classList.toggle("closed");
	await alert(signInFormData.get("Email"));
}
for(let submit of document.querySelectorAll("#submit")){
    submit.addEventListener("click", function(){
        let doc =  window.top.document;
        doc.querySelector("#open-button").hidden = true;
        let table = doc.createElement('table');
        table.setAttribute("class", "table");
        table.setAttribute("id", "table");
        table.insertAdjacentHTML("beforeend",`
        <th>id</th>
        <th>time</th>`);
        for (let i = 0; i < 3; i++){
            table.insertAdjacentHTML("beforeend",`
            <tr>
                <td>${i}</td>
                <td>${Date.now()}</td>
            </tr>`)
        }
        // tableDiv.textContent = "rtyuiuytr";
        // let table = tableDiv.createElement('table');
        // table.setAttribute("class", "table");
        // table.setAttribute("id", "table");
        // let th = doc.createElement("th");
        // th.setAttribute("class", "th");
        // th.setAttribute("id", "th");
        // let thd1 = doc.createElement("td");
        // thd1.textContent = "id";
        // let thd2 = doc.createElement("td");
        // thd2.textContent = "time";
        doc.querySelector("#header").after(table);
        // doc.body.querySelector("#table-div").append(th);
        // doc.body.querySelector("#th").append(thd1);
        // doc.body.querySelector("#th").append(thd2);
        
}
    )};
socket.onmessage = function(event) {
	let message = event.data; 
	// document.querySelector("#table").
	console.log(message);
	  }
