import { type FeedPost, Post } from "@/Post";
import { useCookies } from "react-cookie";
import { useNavigate } from "react-router";
import useSWR, { mutate } from "swr";
// import gohper from "./../public/gohper.svg";
import { CreatePostForm } from "@/CreatePostForm";
import { Button } from "./components/ui/button";

import { API_URL, fetcher } from "@/lib/api";
import { ModeToggle } from "./components/Mode-Toggle";
import { Spinner } from "./components/ui/spinner";

function App() {
  const [cookies, setCookie] = useCookies(["at"]);
  const at = cookies.at;

  let currentUserId = -1;
  if (at) {
    try {
      const payload = JSON.parse(atob(at.split(".")[1]));
      currentUserId = payload.sub;
    } catch (e) {}
  }

  const redirect = useNavigate();

  const { data, error, isLoading } = useSWR<{ data: FeedPost[] }>(
    "/users/feed",
    at ? fetcher(at) : null,
  );

  if (error) return <div>failed to load</div>;
  if (isLoading)
    return (
      <div className="flex h-screen items-center justify-center text-muted-foreground">
        <Spinner />
        <span>Memuat...</span>
      </div>
    );
  if (!data) return <div>Please login first</div>;

  const posts = data.data || [];

  const handleLogout = () => {
    setCookie("at", "", { path: "/" });
    window.location.href = "/";
    return;
  };

  const reFetchData = () => {
    mutate("/users/feed");
  };

  const handleUnfollow = async (userId: number) => {
    try {
      await fetch(`${API_URL}/users/${userId}/unfollow`, {
        method: "PUT",
        headers: {
          Authorization: `Bearer ${at}`,
        },
      });
      mutate("/users/feed");
    } catch (err) {
      console.error("Failed to unfollow user:", err);
    }
  };

  const handleClickPost = (id: number) => () => redirect(`/post/${id}`);

  return (
    <div className="mx-auto max-w-3xl px-4 py-8 md:py-12 flex flex-col gap-8">
      <nav className="flex items-center justify-between border-b pb-4">
        <div className="flex items-center gap-2">
          <h1 className="text-2xl font-bold tracking-tight">SocialNetwork</h1>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" onClick={() => redirect("/explore")}>
            Explore
          </Button>
          <Button variant="outline" onClick={handleLogout}>
            Logout
          </Button>
          <ModeToggle />
        </div>
      </nav>

      <div>
        <h2 className="text-xl font-semibold tracking-tight">Welcome back!</h2>
        <p className="text-muted-foreground mt-1">
          This is a social media platform for nerdy people.
        </p>
      </div>

      <CreatePostForm onFetchPosts={reFetchData} />

      <div className="flex flex-col gap-6">
        {posts.map((post) => (
          <Post
            key={post.id}
            post={post}
            onClick={handleClickPost(post.id)}
            onUnfollow={post.user_id !== currentUserId ? handleUnfollow : undefined}
          />
        ))}

        {posts.length === 0 && (
          <div className="rounded-lg border border-dashed p-12 text-center text-muted-foreground">
            No posts yet, start following someone or post something.
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
