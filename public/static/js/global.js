// global.js script

const thModal = document.getElementById("th-modal");
const thCloseBtn = document.getElementById("close-th-btn");
const thAdd = document.getElementById("add-th-btn");
const replyBtn = document.getElementById("create-post-btn");
const threadCreateBtn = document.getElementById("create-thread-btn");

if (thModal && thCloseBtn && thAdd) {
  // const create thre modal
  thAdd.addEventListener("click", () => {
    thModal.style.display = "flex";
  });

  thCloseBtn.addEventListener("click", (e) => {
    e.preventDefault();

    thModal.style.display = "none";
  });
}

if (threadCreateBtn) {
  // create thread button
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

    // refresh the page after the thread is secessfully created
    if (response.status == 201) {
      window.location.reload();
    }
  };
}
// reply to thread by posting

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

    // refresh the page after the thread is secessfully created
    if (response.status == 201) {
      window.location.reload();
    }
  };
}
