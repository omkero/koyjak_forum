// global.js script

const thModal = document.getElementById("th-modal");
const thCloseBtn = document.getElementById("close-th-btn");
const thAdd = document.getElementById("add-th-btn");
const replyBtn = document.getElementById("create-post-btn");
const threadCreateBtn = document.getElementById("create-thread-btn");

const singupBtn = document.getElementById("signup-btn");
const responseBox = document.getElementById("response-box")
const responseMsg = document.getElementById("resnponse-msg")

if (thModal && thCloseBtn && thAdd) {

  thAdd.addEventListener("click", () => {
    thModal.style.display = "flex";
  });

  thCloseBtn.addEventListener("click", (e) => {
    e.preventDefault();

    thModal.style.display = "none";
  });
}

if (threadCreateBtn) {
  // thread button
  threadCreateBtn.onclick = async function () {
    const threadTitle = document.getElementById("thread-title").value;
    const threadContent = document.getElementById("thread-content").value;

    const response = await fetch("/create_thread", {
      method: "POST",
      headers: {
        "Content-type": "application/json",
      },
      body: JSON.stringify({
        thread_title: threadTitle,
        thread_content: threadContent,
      }),
    });

    const datas = await response.json();
    console.log(datas);

    if (response.status == 201) {
      window.location.reload();
    }
  };
}

if (replyBtn) {
  replyBtn.onclick = async function () {
    const postTitle = document.getElementById("post-title").value;
    const postContent = document.getElementById("post-content").value;

    const response = await fetch("/create_post", {
      method: "POST",
      headers: {
        "Content-type": "application/json",
      },
      body: JSON.stringify({
        post_title: postTitle,
        post_content: postContent,
      }),
    });

    const datas = await response.json();

    if (response.status == 201) {
      window.location.reload();
    }
  };
}

if (singupBtn && responseBox && responseMsg) {
  singupBtn.onclick = async () => {
    try {
      const registerUsernameInput = document.getElementById("reg-username").value;
      const registerEmailInput = document.getElementById("reg-email").value;
      const registerPasswordInput = document.getElementById("reg-password").value;

      if (!registerUsernameInput) {
        responseBox.style.display = "flex";
        responseMsg.innerText = "username is required !!";
        return;

      }
      if (!registerEmailInput) {
        responseBox.style.display = "flex";
        responseMsg.innerText = "email address is required !!";
        return;
      }
      if (!registerPasswordInput) {
        responseBox.style.display = "flex";
        responseMsg.innerText = "password is required !!";
        return;
      }

      responseBox.style.display = "none";
      const response = await fetch("/auth/signup", {
        method: "POST",
        headers: {
          "Content-type": "application/json",
        },
        body: JSON.stringify({
          username: registerUsernameInput,
          email_address: registerEmailInput,
          password: registerPasswordInput
        }),
      })
      const data = await response.json();
      console.log(response.status);
      console.log(data)

      if (response.status !== 201) {
        responseBox.style.display = "flex";
        responseMsg.innerHTML = data?.message;
      }

      if (response.status == 201) {
        window.location.href = "/";
      }


    } catch (err) {

      console.error(err)
    }
  }
}