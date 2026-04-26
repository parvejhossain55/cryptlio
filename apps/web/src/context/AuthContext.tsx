"use client";

import React, { createContext, useContext, useState, useEffect } from "react";
import { useRouter } from "next/navigation";

type UserRole = "user" | "merchant" | "admin" | null;

interface User {
    id: string;
    name: string;
    email: string;
    role: UserRole;
}

interface AuthContextType {
    user: User | null;
    login: (email: string, role: UserRole) => Promise<void>;
    logout: () => void;
    isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const router = useRouter();

    useEffect(() => {
        // Simulate checking for existing session
        const savedUser = localStorage.getItem("cryplio_user");
        if (savedUser) {
            setUser(JSON.parse(savedUser));
        }
        setIsLoading(false);
    }, []);

    const login = async (email: string, role: UserRole) => {
        setIsLoading(true);
        // Simulate API call
        await new Promise((resolve) => setTimeout(resolve, 1500));

        const userData: User = {
            id: "1",
            name: email.split("@")[0],
            email: email,
            role: role || "user",
        };

        setUser(userData);
        localStorage.setItem("cryplio_user", JSON.stringify(userData));
        setIsLoading(false);

        // Redirect based on role
        if (userData.role === "admin") router.push("/admin/dashboard");
        else if (userData.role === "merchant") router.push("/merchant/dashboard");
        else router.push("/user/dashboard");
    };

    const logout = () => {
        setUser(null);
        localStorage.removeItem("cryplio_user");
        router.push("/login");
    };

    return (
        <AuthContext.Provider value={{ user, login, logout, isLoading }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
};
