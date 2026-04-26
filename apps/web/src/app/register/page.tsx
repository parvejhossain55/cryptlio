"use client";

import React, { useState } from "react";
import Link from "next/link";
import { Mail, Lock, User, Eye, EyeOff, Loader2, ArrowRight } from "lucide-react";
import AuthLayout from "@/components/auth/AuthLayout";

const RegisterPage = () => {
    const [showPassword, setShowPassword] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [step, setStep] = useState(1);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        // Simulate registration
        setTimeout(() => {
            setIsLoading(false);
            setStep(2); // Success step
        }, 2000);
    };

    if (step === 2) {
        return (
            <AuthLayout
                title="Verification"
                subtitle="Please check your email to verify your account."
            >
                <div className="py-10">
                    <div className="w-24 h-24 bg-primary/10 rounded-3xl flex items-center justify-center mb-10 shadow-2xl shadow-primary/5">
                        <Mail className="text-primary w-12 h-12" />
                    </div>
                    <p className="text-text-dim text-lg leading-relaxed mb-10">
                        We sent a verification link to your email. Please click the link to continue.
                    </p>
                    <div className="space-y-4">
                        <button
                            type="button"
                            className="w-full bg-white text-background py-5 rounded-2xl text-lg font-black transition-all shadow-2xl active:scale-[0.98]"
                            onClick={() => setStep(1)}
                        >
                            Resend Email
                        </button>
                        <Link
                            href="/login"
                            className="block w-full text-center py-4 text-sm font-bold text-text-dim hover:text-white transition-colors"
                        >
                            Back to Login
                        </Link>
                    </div>
                </div>
            </AuthLayout>
        );
    }

    return (
        <AuthLayout
            title="Register"
            subtitle="Create an account to start trading crypto."
        >
            <form onSubmit={handleSubmit} className="space-y-7">
                <div className="space-y-2">
                    <label className="text-xs font-black text-text-dim uppercase tracking-[0.15em] block px-1">Full Name</label>
                    <div className="relative group">
                        <User className="absolute left-0 top-1/2 -translate-y-1/2 text-text-dim w-5 h-5 group-focus-within:text-primary transition-colors" />
                        <input
                            type="text"
                            required
                            placeholder="Enter your full name"
                            className="w-full bg-transparent border-b border-border py-4 pl-8 pr-4 text-base outline-none focus:border-primary transition-all font-medium"
                        />
                    </div>
                </div>

                <div className="space-y-2">
                    <label className="text-xs font-black text-text-dim uppercase tracking-[0.15em] block px-1">Email Address</label>
                    <div className="relative group">
                        <Mail className="absolute left-0 top-1/2 -translate-y-1/2 text-text-dim w-5 h-5 group-focus-within:text-primary transition-colors" />
                        <input
                            type="email"
                            required
                            placeholder="Enter your email"
                            className="w-full bg-transparent border-b border-border py-4 pl-8 pr-4 text-base outline-none focus:border-primary transition-all font-medium"
                        />
                    </div>
                </div>

                <div className="space-y-2">
                    <label className="text-xs font-black text-text-dim uppercase tracking-[0.15em] block px-1">Password</label>
                    <div className="relative group">
                        <Lock className="absolute left-0 top-1/2 -translate-y-1/2 text-text-dim w-5 h-5 group-focus-within:text-primary transition-colors" />
                        <input
                            type={showPassword ? "text" : "password"}
                            required
                            placeholder="Create a password"
                            className="w-full bg-transparent border-b border-border py-4 pl-8 pr-12 text-base outline-none focus:border-primary transition-all font-medium"
                        />
                        <button
                            type="button"
                            onClick={() => setShowPassword(!showPassword)}
                            className="absolute right-0 top-1/2 -translate-y-1/2 text-text-dim hover:text-white transition-colors"
                        >
                            {showPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
                        </button>
                    </div>
                </div>

                <div className="space-y-4 pt-4">
                    <div className="flex items-start space-x-3 group">
                        <div className="relative flex items-center h-5">
                            <input
                                type="checkbox"
                                id="terms"
                                required
                                className="w-5 h-5 rounded-lg border-border bg-surface text-primary focus:ring-primary/20 accent-primary cursor-pointer"
                            />
                        </div>
                        <label htmlFor="terms" className="text-xs text-text-dim leading-relaxed cursor-pointer select-none group-hover:text-white transition-colors">
                            I agree to the <Link href="/terms" className="text-white font-black hover:underline underline-offset-4">Terms and Conditions</Link>.
                        </label>
                    </div>
                </div>

                <div className="pt-4">
                    <button
                        type="submit"
                        disabled={isLoading}
                        className="w-full bg-white text-background py-5 rounded-2xl text-lg font-black transition-all shadow-2xl shadow-white/5 flex items-center justify-center space-x-2 active:scale-[0.98] disabled:opacity-70 group"
                    >
                        {isLoading ? <Loader2 className="w-6 h-6 animate-spin" /> : (
                            <>
                                <span>Register</span>
                                <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
                            </>
                        )}
                    </button>
                </div>

                <div className="pt-6 text-center border-t border-white/5">
                    <p className="text-xs font-medium text-text-dim">
                        Already have an account? <Link href="/login" className="text-white font-black hover:underline underline-offset-4">Login</Link>
                    </p>
                </div>
            </form>
        </AuthLayout>
    );
};

export default RegisterPage;
