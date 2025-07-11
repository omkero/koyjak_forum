// Type assertions for all elements
const thModal = document.getElementById("th-modal") as HTMLElement | null;
const thCloseBtn = document.getElementById(
  "close-th-btn"
) as HTMLButtonElement | null;
const thAdd = document.getElementById("add-th-btn") as HTMLButtonElement | null;
const replyBtn = document.getElementById(
  "create-post-btn"
) as HTMLButtonElement | null;
const threadCreateBtn = document.getElementById(
  "create-thread-btn"
) as HTMLButtonElement | null;
const thParent = document.getElementById("th-parent") as HTMLElement | null;

const signupBtn = document.getElementById(
  "signup-btn"
) as HTMLButtonElement | null;
const signInBtn = document.getElementById(
  "sginin-btn"
) as HTMLButtonElement | null;

const responseBoxError: any = document.getElementById(
  "response-box-error"
) as HTMLElement | null;
const responseMsgError: any = document.getElementById(
  "resnponse-msg-error"
) as HTMLElement | null;
const responseBoxSuccess: any = document.getElementById(
  "response-box-success"
) as HTMLElement | null;
const responseMsgSuccess: any = document.getElementById(
  "resnponse-msg-success"
) as HTMLElement | null;

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
  threadCreateBtn.onclick = async () => {
    const threadTitleInput = document.getElementById(
      "thread-title"
    ) as HTMLInputElement | null;
    const threadContentInput = document.getElementById(
      "thread-content"
    ) as HTMLTextAreaElement | null;
    const threadForumTitle = document.getElementById(
      "forum-title"
    ) as HTMLTextAreaElement | null;

    if (!threadTitleInput || !threadContentInput || !threadForumTitle) return;

    const threadTitle = threadTitleInput.value;
    const threadContent = threadContentInput.value;
    const forumTitle = threadForumTitle.value;

    const response = await fetch("/create_thread", {
      method: "POST",
      headers: {
        "Content-type": "application/json",
      },
      body: JSON.stringify({
        thread_title: threadTitle,
        thread_content: threadContent,
        forum_title: forumTitle,
      }),
    });

    const datas = await response.json();
    console.log(datas);

    if (response.status === 201) {
      window.location.reload();
    }
  };
}

if (replyBtn && thParent) {
  replyBtn.onclick = async () => {
    const postTitleInput = document.getElementById(
      "post-title"
    ) as HTMLInputElement | null;
    const postContentInput = document.getElementById(
      "post-content"
    ) as HTMLTextAreaElement | null;

    if (!postTitleInput || !postContentInput) return;

    const postTitle = postTitleInput.value;
    const postContent = postContentInput.value;
    const threadToken = thParent.dataset.threadTk;

    const response = await fetch("/create_post", {
      method: "POST",
      headers: {
        "Content-type": "application/json",
      },
      body: JSON.stringify({
        post_title: postTitle,
        post_content: postContent,
        thread_token: threadToken,
      }),
    });

    const datas = await response.json();

    if (response.status === 201) {
      window.location.reload();
    }
  };
}

if (
  signInBtn &&
  responseBoxError &&
  responseMsgError &&
  responseBoxSuccess &&
  responseMsgSuccess
) {
  signInBtn.onclick = async () => {
    try {
      const emailInputHtml = document.getElementById(
        "sginin-email"
      ) as HTMLInputElement | null;
      const passwordInputHtml = document.getElementById(
        "sginin-password"
      ) as HTMLInputElement | null;

      if (!emailInputHtml || !passwordInputHtml) return;

      const emailInput = emailInputHtml.value;
      const passwordInput = passwordInputHtml.value;

      if (!emailInput) {
        responseBoxError.style.display = "flex";
        responseMsgError.innerText = "email is required !!";
        return;
      }

      if (!passwordInput) {
        responseBoxError.style.display = "flex";
        responseMsgError.innerText = "password is required !!";
        return;
      }

      const response = await fetch("/auth/signin", {
        method: "POST",
        headers: {
          "Content-type": "application/json",
        },
        body: JSON.stringify({
          email_address: emailInput,
          password: passwordInput,
        }),
      });

      const data = await response.json();

      if (response.status === 200) {
        window.location.href = "/";
      }
      if (response.status !== 200) {
        responseBoxSuccess.style.display = "none";
        responseBoxError.style.display = "flex";
        responseMsgError.innerText = data?.message;
      }
    } catch (err) {
      // console.error(err);
    }
  };
}

if (
  signupBtn &&
  responseBoxError &&
  responseMsgError &&
  responseBoxSuccess &&
  responseMsgSuccess
) {
  signupBtn.onclick = async () => {
    try {
      const usernameInput = document.getElementById(
        "reg-username"
      ) as HTMLInputElement | null;
      const emailInput = document.getElementById(
        "reg-email"
      ) as HTMLInputElement | null;
      const passwordInput = document.getElementById(
        "reg-password"
      ) as HTMLInputElement | null;

      if (!usernameInput || !emailInput || !passwordInput) return;

      const username = usernameInput.value;
      const email_address = emailInput.value;
      const password = passwordInput.value;

      if (!username) {
        responseBoxError.style.display = "flex";
        responseMsgError.innerText = "username is required !!";
        return;
      }

      if (!email_address) {
        responseBoxError.style.display = "flex";
        responseMsgError.innerText = "email address is required !!";
        return;
      }

      if (!password) {
        responseBoxError.style.display = "flex";
        responseMsgError.innerText = "password is required !!";
        return;
      }

      const response = await fetch("/auth/signup", {
        method: "POST",
        headers: {
          "Content-type": "application/json",
        },
        body: JSON.stringify({ username, email_address, password }),
      });

      const data = await response.json();

      if (response.status !== 201) {
        responseBoxSuccess.style.display = "none";
        responseBoxError.style.display = "flex";
        responseMsgError.innerText = data?.message;
      } else {
        responseBoxError.style.display = "none";
        responseBoxSuccess.style.display = "flex";
        responseMsgSuccess.innerText =
          "You have successfully created an account. You can login now.";
      }
    } catch (err) {
      //console.error(err);
    }
  };
}

function CloseErrAlert() {
  if (responseBoxError) responseBoxError.style.display = "none";
}

function CloseSuccAlert() {
  if (responseBoxSuccess) responseBoxSuccess.style.display = "none";
}
