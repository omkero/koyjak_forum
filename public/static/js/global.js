// global.js script

const thModal = document.getElementById("th-modal");
const thCloseBtn = document.getElementById("close-th-btn");
const thAdd = document.getElementById("add-th-btn");
const replyBtn = document.getElementById("create-post-btn");
const threadCreateBtn = document.getElementById("create-thread-btn");

const singupBtn = document.getElementById("signup-btn");

const responseBoxError = document.getElementById("response-box-error")
const responseMsgError = document.getElementById("resnponse-msg-error")
const responseBoxSuccess = document.getElementById("response-box-success")
const responseMsgSuccess = document.getElementById("resnponse-msg-success")

const signInBtn = document.getElementById("sginin-btn");


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

if (singupBtn && responseBoxError && responseMsgError && responseBoxSuccess && responseMsgSuccess) {
  singupBtn.onclick = async () => {
    try {

      const registerUsernameInput = document.getElementById("reg-username").value;
      const registerEmailInput = document.getElementById("reg-email").value;
      const registerPasswordInput = document.getElementById("reg-password").value;

      if (!registerUsernameInput) {
        responseBoxError.style.display = "flex";
        responseMsgError.innerText = "username is required !!";
        return;

      }

      if (!registerEmailInput) {
        responseBoxError.style.display = "flex";
        responseMsgError.innerText = "email address is required !!";
        return;
      }

      if (!registerPasswordInput) {
        responseBoxError.style.display = "flex";
        responseMsgError.innerText = "password is required !!";
        return;
      }

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

      if (response.status !== 201) {
        responseBoxSuccess.style.display = "none";
        responseBoxError.style.display = "flex";
        responseMsgError.innerText = data?.message;
      }

      if (response.status == 201) {
        responseBoxError.style.display = "none";
        responseBoxSuccess.style.display = "flex";
        responseMsgSuccess.innerText = "You have successfully created an account you can login now."
      }
    } catch (err) {
      console.log(err)
    }
  }
}

function CloseErrAlert() {
  responseBoxError.style.display = "none";
}

function CloseSuccAlert() {
  responseBoxSuccess.style.display = "none";
}

if (signInBtn) {
    signInBtn.onclick = async () => {
    try {

      const signinEmailInput = document.getElementById("sginin-email").value;
      const signinPasswordInput = document.getElementById("sginin-password").value;

      const response = await fetch("/auth/signin", {
        method: "POST",
        headers: {
          "Content-type": "application/json",
        },
        body: JSON.stringify({
          email_address: signinEmailInput,
          password: signinPasswordInput
        }),
      })

      const data = await response.json();
      console.log(data);

      if (response.status == 200) {
        window.location.href = "/"
      }
    } catch (err) {
      console.log(err)
    }
  }
}