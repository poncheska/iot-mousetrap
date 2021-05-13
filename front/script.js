// let socket = new WebSocket("ws://smart-mousetrap.herokuapp.com/mousetraps/ws");
let modalSignIn = document.querySelector("#modal-sign-in"),
modalJoin = document.querySelector("#modal-join"),
modalOverlay = document.querySelector("#modal-overlay"),
closeButtonSignIn = document.querySelector("#close-button-sign-in"),
closeButtonJoin = document.querySelector("#close-button-join"),
openButtonSignIn = document.querySelector("#open-button-sign-in");
openButtonJoin = document.querySelector("#open-button-join");
async function join() {
    let formdata = new FormData( document.forms.join);
    let response = await fetch("/org/sign-up", {
        method: 'POST',
        headers: {'Content-Type': 'application/json;charset=utf-8'
    },
    body: JSON.stringify({"name":formdata.get("Email"), "pass":formdata.get("Password")})
    });
    if (!response.ok){
        let errorMessage = document.createElement('div');
        errorMessage.className = "error";
        errorMessage.innerHTML = `<strong>${JSON.parse(response).message}</strong>`
        document.querySelector("#form-name").after(errorMessage);
    }
}
// document.forms.form
async function signIn(form) {
    if (form.id == "join-form"){
        await join();
    }
    let formdata = new FormData(form);
    let response = fetch("/org/sign-in", {
        method: 'POST',
        headers: {'Content-Type': 'application/json;charset=utf-8'
    },
    body: JSON.stringify({"name":formdata.get("Email"), "pass":formdata.get("Password")})
    });
    if (response.ok){
        localStorage.setItem('token', JSON.parse(response).token);
        modalSignIn.classList.toggle("closed");
        modalOverlay.classList.toggle("closed");
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
        doc.querySelector("#header").after(table)
    }
    else{
        let errorMessage = document.createElement('div');
        errorMessage.className = "error";
        errorMessage.innerHTML = `<strong>${JSON.parse(response).message}</strong>`
        document.querySelector("#form-name").after(errorMessage);
    }
    return false
}
function buttonDefinition(modal, modalOverlay, closeButton, openButton) {
	closeButton.addEventListener("click", function() {
		modal.classList.toggle("closed");
		modalOverlay.classList.toggle("closed");
        //modal.querySelector("#Email").value = "";
        //modal.querySelector("#Password").value = "";
	});

    openButton.addEventListener("click", function() {
        modal.classList.toggle("closed");
        modalOverlay.classList.toggle("closed");
    });
    }
    buttonDefinition(modalSignIn, modalOverlay, closeButtonSignIn, openButtonSignIn);
    buttonDefinition(modalJoin, modalOverlay, closeButtonJoin, openButtonJoin);
// socket.onmessage = function(event) {
// 	let message = event.data; 
// 	// document.querySelector("#table").
// 	console.log(message);
// 	  }
document.forms.join.addEventListener("submit", function() {signIn(document.forms.join)});
document.forms.signIn.addEventListener("submit", function() {signIn(document.forms.signIn)});
