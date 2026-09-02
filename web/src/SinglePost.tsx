import { fetcher } from "@/lib/api";
import { type FeedPost } from "@/Post";
import { useCookies } from "react-cookie";
import { useNavigate, useParams } from "react-router";
import useSWR from "swr";
import { Button } from "./components/ui/button";
import { Separator } from "./components/ui/separator";
import { Spinner } from "./components/ui/spinner";

export const SinglePost = () => {
  const { postID } = useParams();
  const [cookies] = useCookies(["at"]);
  const at = cookies.at;

  const redirect = useNavigate();

  const { data, error, isLoading } = useSWR<{ data: FeedPost }>(
    "/posts/" + postID,
    at ? fetcher(at) : null,
  );

  if (error)
    return (
      <div className="flex h-screen items-center justify-center text-muted-foreground">
        Failed to load post.
      </div>
    );

  if (isLoading)
    return (
      <div className="flex h-screen items-center justify-center text-muted-foreground">
        <Spinner />
        <span>Memuat...</span>
      </div>
    );

  if (!data)
    return (
      <div className="flex h-screen items-center justify-center text-muted-foreground">
        Please login first.
      </div>
    );

  const { data: post } = data;

  const formatDate = (dateStr: string) => {
    const d = new Date(dateStr);
    return d.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  };

  const formatTime = (dateStr: string) => {
    const d = new Date(dateStr);
    return d.toLocaleTimeString("en-US", {
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  return (
    <div className="mx-auto max-w-3xl px-4 py-8 md:py-12">
      {/* Back button */}
      <Button
        variant="ghost"
        className="mb-6 gap-1.5 px-2 text-muted-foreground hover:text-foreground"
        onClick={() => redirect("/")}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <path d="m15 18-6-6 6-6" />
        </svg>
        Back to feed
      </Button>

      {/* Post header */}
      <article>
        <header className="mb-6">
          <h1 className="text-3xl font-bold tracking-tight">
            {post.title || "Untitled"}
          </h1>
          <p className="mt-2 text-sm text-muted-foreground">
            {formatDate(post.created_at)} · {formatTime(post.created_at)}
          </p>

          {post.tags && post.tags.length > 0 && (
            <div className="mt-3 flex flex-wrap gap-1.5">
              {post.tags.map((tag) => (
                <span
                  key={tag}
                  className="rounded-full bg-secondary px-2.5 py-0.5 text-xs font-medium text-secondary-foreground"
                >
                  {tag}
                </span>
              ))}
            </div>
          )}
        </header>

        {/* Post content */}
        <div className="prose prose-neutral max-w-none leading-relaxed text-foreground/90">
          <p className="whitespace-pre-wrap">{post.content}</p>
        </div>
      </article>

      <Separator className="my-8" />

      {/* Comments section */}
      <section>
        <h2 className="mb-4 text-lg font-semibold tracking-tight">
          Comments
          {post.comments && post.comments.length > 0 && (
            <span className="ml-1.5 text-sm font-normal text-muted-foreground">
              ({post.comments.length})
            </span>
          )}
        </h2>

        {!post.comments || post.comments.length === 0 ? (
          <div className="rounded-lg border border-dashed p-8 text-center text-sm text-muted-foreground">
            No comments yet.
          </div>
        ) : (
          <div className="flex flex-col gap-4">
            {post.comments.map((comment) => (
              <div key={comment.id} className="rounded-lg border bg-card p-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium">
                    {comment.user?.username || "Anonymous"}
                  </span>
                  <span className="text-xs text-muted-foreground">
                    {formatDate(comment.created_at)}
                  </span>
                </div>
                <p className="mt-1.5 text-sm leading-relaxed text-foreground/80">
                  {comment.content}
                </p>
              </div>
            ))}
          </div>
        )}
      </section>
    </div>
  );
};
