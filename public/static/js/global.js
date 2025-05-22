// global.js script

const thModal = document.getElementById("th-modal") 
const thCloseBtn = document.getElementById("close-th-btn")
const thAdd = document.getElementById("add-th-btn")

// const create thre modal 

const threadCreateBtn =  document.getElementById("create-thread-btn")

thAdd.addEventListener("click", () => {
    thModal.style.display = "flex"
})

thCloseBtn.addEventListener("click", () => {
    thModal.style.display = "none"
})

// create thread button
threadCreateBtn.onclick = async function () {
    const threadTitle = document.getElementById("thread-title").value
    const threadContent = document.getElementById("thread-content").value

    const response = await fetch("/create_thread", {
        method: "POST",
        headers: {
            "Content-type": "application/json"
        },
        body: JSON.stringify({
            thread_title: threadTitle,
            thread_content: threadContent
        })
    })

    const datas = await response.json();
    console.log(datas)

    // refresh the page after the thread is secessfully created
    if (response.status == 201) {
        window.location.reload()
    }
}