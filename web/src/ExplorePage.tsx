import { type FeedPost } from "@/Post";
import { API_URL, fetcher } from "@/lib/api";
import { useCookies } from "react-cookie";
import { useNavigate } from "react-router";
import useSWR, { mutate } from "swr";
import { Button } from "./components/ui/button";
import { Spinner } from "./components/ui/spinner";
import { useState } from "react";

interface ExplorePostCardProps {
  post: FeedPost;
  onFollow: (userId: number) => void;
  isFollowing: boolean;
  onClick: () => void;
}

const ExplorePostCard: React.FC<ExplorePostCardProps> = ({
  post,
  onFollow,
  isFollowing,
  onClick,
}) => {
  const date = new Date(post.created_at);
  const timeAgo = getTimeAgo(date);

  return (
    <article className="group rounded-xl border bg-card p-5 transition-all duration-200 hover:border-foreground/20 hover:shadow-sm">
      {/* Header: user info + follow button */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <div className="flex h-8 w-8 items-center justify-center rounded-full bg-secondary text-xs font-semibold text-secondary-foreground">
            {(post.user?.username ?? "?")[0].toUpperCase()}
          </div>
          <div>
            <span className="font-medium text-foreground">
              {post.user?.username ?? "Unknown"}
            </span>
            <span className="ml-1.5">· {timeAgo}</span>
          </div>
        </div>
        <Button
          variant="outline"
          size="sm"
          disabled={isFollowing}
          onClick={(e) => {
            e.stopPropagation();
            onFollow(post.user_id);
          }}
          className="text-xs"
        >
          {isFollowing ? (
            <span className="flex items-center gap-1.5">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="14"
                height="14"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <polyline points="20 6 9 17 4 12" />
              </svg>
              Following
            </span>
          ) : (
            <span className="flex items-center gap-1.5">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="14"
                height="14"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
                <circle cx="9" cy="7" r="4" />
                <line x1="19" y1="8" x2="19" y2="14" />
                <line x1="22" y1="11" x2="16" y2="11" />
              </svg>
              Follow
            </span>
          )}
        </Button>
      </div>

      {/* Post content - clickable */}
      <div className="mt-3 cursor-pointer" onClick={onClick}>
        <h2 className="text-base font-semibold leading-snug tracking-tight group-hover:text-foreground/80">
          {post.title || "Untitled"}
        </h2>
        <p className="mt-1 line-clamp-3 text-sm leading-relaxed text-muted-foreground">
          {post.content}
        </p>
      </div>

      {/* Tags */}
      {post.tags && post.tags.length > 0 && (
        <div className="mt-3 flex flex-wrap gap-1.5">
          {post.tags.map((tag) => (
            <span
              key={tag}
              className="rounded-full bg-secondary px-2 py-0.5 text-xs font-medium text-secondary-foreground"
            >
              {tag}
            </span>
          ))}
        </div>
      )}

      {/* Footer */}
      <div className="mt-3 flex items-center gap-1.5 text-xs text-muted-foreground">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="14"
          height="14"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
        </svg>
        <span>
          {post.comments_count}{" "}
          {post.comments_count === 1 ? "comment" : "comments"}
        </span>
      </div>
    </article>
  );
};

export default function ExplorePage() {
  const [cookies] = useCookies(["at"]);
  const at = cookies.at;
  const redirect = useNavigate();
  const [followedUsers, setFollowedUsers] = useState<Set<number>>(new Set());

  const { data, error, isLoading } = useSWR<{ data: FeedPost[] }>(
    "/users/explore",
    at ? fetcher(at) : null,
  );

  if (error)
    return (
      <div className="flex h-screen items-center justify-center text-muted-foreground">
        Failed to load.
      </div>
    );

  if (isLoading)
    return (
      <div className="flex h-screen items-center justify-center text-muted-foreground">
        <Spinner />
        <span className="ml-2">Memuat...</span>
      </div>
    );

  if (!data)
    return (
      <div className="flex h-screen items-center justify-center text-muted-foreground">
        Please login first.
      </div>
    );

  const posts = data.data || [];

  const handleFollow = async (userId: number) => {
    try {
      await fetch(`${API_URL}/users/${userId}/follow`, {
        method: "PUT",
        headers: {
          Authorization: `Bearer ${at}`,
        },
      });
      setFollowedUsers((prev) => new Set(prev).add(userId));
      // Refresh both feeds
      mutate("/users/explore");
      mutate("/users/feed");
    } catch (err) {
      console.error("Failed to follow user:", err);
    }
  };

  return (
    <div className="mx-auto max-w-3xl px-4 py-8 md:py-12 flex flex-col gap-8">
      {/* Back nav */}
      <div className="flex items-center justify-between border-b pb-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Explore</h1>
          <p className="text-sm text-muted-foreground mt-0.5">
            Discover posts from people you don't follow yet.
          </p>
        </div>
        <Button variant="outline" onClick={() => redirect("/")}>
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
            className="mr-1.5"
          >
            <path d="m15 18-6-6 6-6" />
          </svg>
          Back to feed
        </Button>
      </div>

      {/* Posts */}
      <div className="flex flex-col gap-6">
        {posts.map((post) => (
          <ExplorePostCard
            key={post.id}
            post={post}
            onFollow={handleFollow}
            isFollowing={followedUsers.has(post.user_id)}
            onClick={() => redirect(`/post/${post.id}`)}
          />
        ))}

        {posts.length === 0 && (
          <div className="rounded-lg border border-dashed p-12 text-center text-muted-foreground">
            No new users to discover. You're following everyone!
          </div>
        )}
      </div>
    </div>
  );
}

function getTimeAgo(date: Date): string {
  const now = new Date();
  const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);

  if (seconds < 60) return "just now";

  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;

  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;

  const days = Math.floor(hours / 24);
  if (days < 7) return `${days}d ago`;

  return date.toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
  });
}
