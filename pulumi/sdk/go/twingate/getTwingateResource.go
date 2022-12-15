// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package twingate

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func LookupTwingateResource(ctx *pulumi.Context, args *LookupTwingateResourceArgs, opts ...pulumi.InvokeOption) (*LookupTwingateResourceResult, error) {
	opts = pkgInvokeDefaultOpts(opts)
	var rv LookupTwingateResourceResult
	err := ctx.Invoke("twingate:index/getTwingateResource:getTwingateResource", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getTwingateResource.
type LookupTwingateResourceArgs struct {
	Id        string                        `pulumi:"id"`
	Protocols []GetTwingateResourceProtocol `pulumi:"protocols"`
}

// A collection of values returned by getTwingateResource.
type LookupTwingateResourceResult struct {
	Address         string                        `pulumi:"address"`
	Id              string                        `pulumi:"id"`
	Name            string                        `pulumi:"name"`
	Protocols       []GetTwingateResourceProtocol `pulumi:"protocols"`
	RemoteNetworkId string                        `pulumi:"remoteNetworkId"`
}

func LookupTwingateResourceOutput(ctx *pulumi.Context, args LookupTwingateResourceOutputArgs, opts ...pulumi.InvokeOption) LookupTwingateResourceResultOutput {
	return pulumi.ToOutputWithContext(context.Background(), args).
		ApplyT(func(v interface{}) (LookupTwingateResourceResult, error) {
			args := v.(LookupTwingateResourceArgs)
			r, err := LookupTwingateResource(ctx, &args, opts...)
			var s LookupTwingateResourceResult
			if r != nil {
				s = *r
			}
			return s, err
		}).(LookupTwingateResourceResultOutput)
}

// A collection of arguments for invoking getTwingateResource.
type LookupTwingateResourceOutputArgs struct {
	Id        pulumi.StringInput                    `pulumi:"id"`
	Protocols GetTwingateResourceProtocolArrayInput `pulumi:"protocols"`
}

func (LookupTwingateResourceOutputArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupTwingateResourceArgs)(nil)).Elem()
}

// A collection of values returned by getTwingateResource.
type LookupTwingateResourceResultOutput struct{ *pulumi.OutputState }

func (LookupTwingateResourceResultOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupTwingateResourceResult)(nil)).Elem()
}

func (o LookupTwingateResourceResultOutput) ToLookupTwingateResourceResultOutput() LookupTwingateResourceResultOutput {
	return o
}

func (o LookupTwingateResourceResultOutput) ToLookupTwingateResourceResultOutputWithContext(ctx context.Context) LookupTwingateResourceResultOutput {
	return o
}

func (o LookupTwingateResourceResultOutput) Address() pulumi.StringOutput {
	return o.ApplyT(func(v LookupTwingateResourceResult) string { return v.Address }).(pulumi.StringOutput)
}

func (o LookupTwingateResourceResultOutput) Id() pulumi.StringOutput {
	return o.ApplyT(func(v LookupTwingateResourceResult) string { return v.Id }).(pulumi.StringOutput)
}

func (o LookupTwingateResourceResultOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v LookupTwingateResourceResult) string { return v.Name }).(pulumi.StringOutput)
}

func (o LookupTwingateResourceResultOutput) Protocols() GetTwingateResourceProtocolArrayOutput {
	return o.ApplyT(func(v LookupTwingateResourceResult) []GetTwingateResourceProtocol { return v.Protocols }).(GetTwingateResourceProtocolArrayOutput)
}

func (o LookupTwingateResourceResultOutput) RemoteNetworkId() pulumi.StringOutput {
	return o.ApplyT(func(v LookupTwingateResourceResult) string { return v.RemoteNetworkId }).(pulumi.StringOutput)
}

func init() {
	pulumi.RegisterOutputType(LookupTwingateResourceResultOutput{})
}