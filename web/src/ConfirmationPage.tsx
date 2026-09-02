import { useNavigate, useParams } from "react-router";
import { API_URL } from "@/lib/api";
import { Button } from "./components/ui/button";
import { cn } from "./lib/utils";

export const ConfirmationPage = () => {
  const { token = "" } = useParams();
  const redirect = useNavigate();

  const handleConfirm = async () => {
    const response = await fetch(`${API_URL}/users/activate/${token}`, {
      method: "PUT",
    });

    if (response.ok) {
      redirect("/");
    } else {
      alert("Failed to confirm token");
    }
  };

  return (
    <div className="relative w-full overflow-hidden md:h-screen">
      <div
        className={cn(
          "relative mx-auto flex min-h-screen w-full max-w-md flex-col justify-center p-6 md:p-8",
        )}
      >
        <div className="fade-in slide-in-from-bottom-4 w-full animate-in space-y-4 duration-600">
          <div className="flex flex-col space-y-1">
            <h1 className="font-bold text-2xl">Konfirmasi Akun</h1>
            <p className="text-base text-muted-foreground">
              Konfirmasi akun anda untuk dapat menggunakan layanan ini.
            </p>
          </div>
          <div className="space-y-2">
            <Button
              onClick={handleConfirm}
              className="w-full cursor-pointer"
              type="button"
            >
              konfirmasi Akun
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
};
